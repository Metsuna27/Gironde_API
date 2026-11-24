const loginForm = document.getElementById("loginForm");
const registerForm = document.getElementById("registerForm");
const message = document.getElementById("message");

function showForm(type) {

    loginForm.hidden = type !== "login";
    registerForm.hidden = type !== "register";

    message.textContent = "";
}

async function sendRequest(endpoint, data) {

    try {

        const response = await fetch(
            `/api/auth/${endpoint}`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            }
        );

        const result = await response.json();

        message.textContent = result.message;

        message.style.color = result.success
            ? "#4ade80"
            : "#f87171";

        return result;

    } catch (error) {

        message.textContent = "Impossible de joindre le serveur";
        message.style.color = "#f87171";

        return null;
    }
}

registerForm.addEventListener("submit", async (event) => {

    event.preventDefault();

    const data = Object.fromEntries(
        new FormData(registerForm)
    );

    const result = await sendRequest("register", data);

    if (result?.success) {
        registerForm.reset();
        showForm("login");

        message.textContent =
            "Compte créé ! Vous pouvez vous connecter.";
    }
});

loginForm.addEventListener("submit", async (event) => {

    event.preventDefault();

    const data = Object.fromEntries(
        new FormData(loginForm)
    );

    const result = await sendRequest("login", data);

    if (result?.success) {
        loginForm.reset();
    }
});