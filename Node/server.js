const express = require("express");
const path = require("path");

const authRoutes = require("./routes/authRoutes");

const app = express();
const PORT = 3000;

app.use(express.json({ limit: "16kb" }));

app.use(express.static(
    path.join(__dirname, "public")
));

app.use("/api/auth", authRoutes);

app.listen(PORT, "127.0.0.1", () => {
    console.log(
        `Serveur Node.js : http://127.0.0.1:${PORT}`
    );
});