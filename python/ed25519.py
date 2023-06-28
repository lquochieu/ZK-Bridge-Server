import hashlib
import json
import base64

def sha512(s):
    return hashlib.sha512(s).digest()

# Base field Z_p
p = 2**255 - 19

def modp_inv(x):
    return pow(x, p-2, p)

# Curve constant
d = -121665 * modp_inv(121666) % p

# Group order
q = 2**252 + 27742317777372353535851937790883648493

def sha512_modq(s):
    return int.from_bytes(sha512(s), "little") % q

## Then follows functions to perform point operations.

# Points are represented as tuples (X, Y, Z, T) of extended
# coordinates, with x = X/Z, y = Y/Z, x*y = T/Z

def point_add(P, Q):
    A, B = (P[1]-P[0]) * (Q[1]-Q[0]) % p, (P[1]+P[0]) * (Q[1]+Q[0]) % p
    C, D = 2 * P[3] * Q[3] * d % p, 2 * P[2] * Q[2] % p
    E, F, G, H = B-A, D-C, D+C, B+A
    return (E*F, G*H, F*G, E*H)

# Computes Q = s * Q
def point_mul(s, P):
    Q = (0, 1, 1, 0)  # Neutral element
    while s > 0:
        if s & 1:
            Q = point_add(Q, P)
        P = point_add(P, P)
        s >>= 1
    return Q

def point_equal(P, Q):
    # x1 / z1 == x2 / z2  <==>  x1 * z2 == x2 * z1
    if (P[0] * Q[2] - Q[0] * P[2]) % p != 0:
        return False
    if (P[1] * Q[2] - Q[1] * P[2]) % p != 0:
        return False
    return True

## Now follows functions for point compression.

# Square root of -1
modp_sqrt_m1 = pow(2, (p-1) // 4, p)

# Compute corresponding x-coordinate, with low bit corresponding to
# sign, or return None on failure
def recover_x(y, sign):
    if y >= p:
        return None
    x2 = (y*y-1) * modp_inv(d*y*y+1)
    if x2 == 0:
        if sign:
            return None
        else:
            return 0

    # Compute square root of x2
    x = pow(x2, (p+3) // 8, p)
    if (x*x - x2) % p != 0:
        x = x * modp_sqrt_m1 % p
    if (x*x - x2) % p != 0:
        return None

    if (x & 1) != sign:
        x = p - x
    return x

g_y = 4 * modp_inv(5) % p
g_x = recover_x(g_y, 0)
G = (g_x, g_y, 1, g_x * g_y % p)

def point_compress(P):
    zinv = modp_inv(P[2])
    x = P[0] * zinv % p
    y = P[1] * zinv % p
    return int.to_bytes(y | ((x & 1) << 255), 32, "little")

def point_decompress(s):
    if len(s) != 32:
        raise Exception("Invalid input length for decompression")
    y = int.from_bytes(s, "little")
    sign = y >> 255
    y &= (1 << 255) - 1

    x = recover_x(y, sign)
    if x is None:
        return None
    else:
        return (x, y, 1, x*y % p)

## These are functions for manipulating the private key.

def secret_expand(secret):
    if len(secret) != 32:
        raise Exception("Bad size of private key")
    h = sha512(secret)
    a = int.from_bytes(h[:32], "little")
    a &= (1 << 254) - 8
    a |= (1 << 254)
    return (a, h[32:])

def secret_to_public(secret):
    (a, dummy) = secret_expand(secret)
    return point_compress(point_mul(a, G))

## The signature function works as below.

def sign(secret, msg):
    a, prefix = secret_expand(secret)
    A = point_compress(point_mul(a, G))
    r = sha512_modq(prefix + msg)
    R = point_mul(r, G)
    Rs = point_compress(R)
    h = sha512_modq(Rs + A + msg)
    s = (r + h * a) % q
    print("pubKey", A.hex())
    print(verify(A, msg, Rs + int.to_bytes(s, 32, "little")))
    return Rs + int.to_bytes(s, 32, "little")

## And finally the verification function.

def verify(public, msg, signature):
    if len(public) != 32:
        raise Exception("Bad public key length")
    if len(signature) != 64:
        Exception("Bad signature length")
    A = point_decompress(public)
    if not A:
        return False
    Rs = signature[:32]
    R = point_decompress(Rs)
    if not R:
        return False
    # A = (19609600535639426967582330360073854330664420980290928614443703354937550235772, 4819101209465356224883271557990864103528016550052741516590013689083114432765, 1, 51169037833951159944323311113976941583233159633603343645972292487679672161224)
    # R = (37880392645989658068752609291251083631204709745327980624473693943571704213899, 4441405635204378107364157117381717143569046592695380934007350497764283041565, 1, 29423111549617446375779108146631000784035658644093037228114776583127667239194)
    s = int.from_bytes(signature[32:], "little")
    print("A", type(A))
    print("R", Rs)
    # print(Rs.hex())
    print(signature[32:].hex())
    if s >= q: return False
    h = sha512_modq(Rs + public + msg)
    sB = point_mul(s, G)
    hA = point_mul(h, A)
    return point_equal(sB, point_add(R, hA))

def convert_base64_to_hex(json_file_path):
    # Read the JSON file
    with open(json_file_path) as file:
        data = json.load(file)

    # Convert base64 fields to hexadecimal strings
    data['dataHash'] = base64.b64decode(data['dataHash']).hex()
    data['parrentSiblings'] = [base64.b64decode(sibling).hex() for sibling in data['parrentSiblings']]
    data['blockHash'] = base64.b64decode(data['blockHash']).hex()
    data['partsHash'] = base64.b64decode(data['partsHash']).hex()
    data['pubKeys'] = base64.b64decode(data['pubKeys']).hex()
    data['R8'] = base64.b64decode(data['R8']).hex()
    data['S'] = base64.b64decode(data['S']).hex()

    return data

def save_fields_to_json(data, file_path):
    """
    Save the provided fields data to a JSON file.
    
    Args:
        data (dict): A dictionary containing the field values.
        file_path (str): The path and filename for the JSON file.
    """
    with open(file_path, "w") as json_file:
        json.dump(data, json_file)

    print("JSON file saved successfully.")
    
data = convert_base64_to_hex("../resources/cosmosHeader/input_go.json")

pubKeys = bytes.fromhex(data['pubKeys'])
R8 = bytes.fromhex(data['R8'])
S = bytes.fromhex(data['S'])

data['A'] = point_decompress(pubKeys)
data['R'] = point_decompress(R8)

save_fields_to_json(data, "../resources/cosmosHeader/input_python.json")
