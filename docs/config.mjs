const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"

export default {
  url: isProduction ? "https://devan.gg" : "http://localhost:4321",
  basePath: isProduction ? "/bookmark" : "/",
  github: "https://github.com/imdevan/bookmark/",
  githubDocs: "https://github.com/imdevan/bookmark/",
  title: "bookmark",
  description: "A bookmark manager for your favorite shell",
}
