const express = require("express");
const client = require("../grpc/authClient");

const router = express.Router();

// Inscription
router.post("/register", (req, res) => {

    const { firstname, lastname, email, password } = req.body;

    if (![firstname, lastname, email, password]
        .every(v => typeof v === "string" && v.length > 0)) {
        return res.status(400).json({
            success: false,
            message: "Tous les champs sont obligatoires"
        });
    }

    client.Register(
        { firstname, lastname, email, password },
        { deadline: Date.now() + 10000 },
        (err, response) => {

            if (err) {
                return res.status(503).json({
                    success: false,
                    message: "Service indisponible"
                });
            }

            return res
                .status(response.success ? 201 : 400)
                .json(response);
        }
    );
});

// Connexion
router.post("/login", (req, res) => {

    const { email, password } = req.body;

    if (![email, password]
        .every(v => typeof v === "string" && v.length > 0)) {
        return res.status(400).json({
            success: false,
            message: "Email et mot de passe obligatoires"
        });
    }

    client.Login(
        { email, password },
        { deadline: Date.now() + 10000 },
        (err, response) => {

            if (err) {
                return res.status(503).json({
                    success: false,
                    message: "Service indisponible"
                });
            }

            return res
                .status(response.success ? 200 : 401)
                .json(response);
        }
    );
});

module.exports = router;