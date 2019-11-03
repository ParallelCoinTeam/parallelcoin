import sirv from "sirv";
import express from "express";
import compression from "compression";
import * as sapper from "@sapper/server";

import "./ass/css/tailwind.css";

const { PORT, NODE_ENV } = process.env;
const dev = NODE_ENV === "development";

const app = express() // You can also use Express
  .use(
    compression({ threshold: 0 }),
    sirv("src/ass", { dev }),
    sapper.middleware()
  );

export default app;

if (!process.env.NOW_REGION) {
  app.listen(PORT, err => {
    if (err) console.log("error", err);
  });
}
