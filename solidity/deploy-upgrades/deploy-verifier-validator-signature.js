const { deployVerifierValidatorSignature } = require("./deploy");

const main = async () => {
    await deployVerifierValidatorSignature();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });


