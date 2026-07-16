const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"

export default {
  url: isProduction ? "https://devan.gg" : "http://localhost:4321",
  basePath: isProduction ? "/aliases" : "/",
  github: "https://github.com/imdevan/aliases/",
  githubDocs: "https://github.com/imdevan/aliases/",
  title: "aliases",
  description: "An alias manager for your favorite shell",
}
