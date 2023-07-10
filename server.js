const express = require("express");
require("dotenv").config();
const compression = require("compression");
const cors = require("cors");
const winston = require("winston");
const morgan = require("morgan");
const mongoose = require("mongoose");
const swaggerUi = require('swagger-ui-express');
const cron = require("node-cron");

const swaggerDocument = require("./src/swagger/swagger.json");
const { cronjobUpdate } = require("./src/controller/Controller");

const app = express();

const logger = winston.createLogger({
    level: 'debug',
    levels: winston.config.npm.levels,
    format: winston.format.combine(
        winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss:ms' }),
        winston.format.colorize({ all: true }),
        winston.format.printf(
            (info) => `${info.timestamp} ${info.level}: ${info.message}`,
        ),
    ),
    transports: [new winston.transports.Console()]
})
const myStream = {
    write: (text) => {
      logger.info(text);
    }
}

app.set("port", process.env.PORT ?? 8000);
app.use(express.json());
app.use(express.urlencoded({extended: true}));
app.use(compression());
app.use(cors());
app.use(morgan("dev", { stream: myStream }));
app.use("/api-docs", swaggerUi.serve, swaggerUi.setup(swaggerDocument));
app.get("/", (req, res) => {
    res.json({ message: "Server is online." });
});

// auto run query and update root
const job = cron.schedule(" * */2 * * * *", async() => { // runs every 2 minutes
    await cronjobUpdate();
}, {
    scheduled: false
});

// job run after 2 minutes
// setTimeout(() => {
//     job.start();
// }, 1000 * 60 * 2);

const MONGODB_URL = process.env.MONGODB_URL ?? 'mongodb://localhost:27017/ZK_BRIDGE';

mongoose.connect(MONGODB_URL, { 
	useNewUrlParser: true,
	useUnifiedTopology: true ,
	socketTimeoutMS: 30000,
    keepAlive: true
}).then(() => {
	//don't show the log when it is test
	if(process.env.NODE_ENV !== "test") {
		logger.info("Connected to %s", MONGODB_URL);
		logger.info("App is running ... \n");
		logger.info("Press CTRL + C to stop the process. \n");
	}

	app.use("/api/", require('./src/router'));

	const PORT = process.env.PORT || 8000;
	app.listen(PORT, () => {
		logger.info(`Server is running on port ${PORT}.`);
	});
}).catch(err => {
    logger.info("App starting error:", err.message);
    process.exit(1);
});