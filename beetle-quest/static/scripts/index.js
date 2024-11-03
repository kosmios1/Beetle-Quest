const session_cookie_name = "my-session";

function sessionIsActive() {
    return localStorage.getItem(session_cookie_name) !== null;
}

function adjustVisibilityBasedOnSession() {
    if (sessionIsActive()) {
        document.getElementById("auth-container").style.display = "none";
    }
}

function logout() {
    localStorage.removeItem(session_cookie_name);
}

document.addEventListener("DOMContentLoaded", adjustVisibilityBasedOnSession);
