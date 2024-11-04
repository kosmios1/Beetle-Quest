const session_cookie_name = "my-session";
const api_endpoint = "/api/v1";

function checkSessionActive() {
    fetch(api_endpoint + "/check_session", {
        method: "GET",
        credentials: "include",
    }).then((response) => {
        adjustVisibilityBasedOnSession(response.ok);
    });
}

function adjustVisibilityBasedOnSession(isSessionActive) {
    if (isSessionActive) {
        document.getElementById("auth-container").style.display = "none";
    } else {
        document.getElementById("auth-container").style.display = "flex";
    }
}

function logout() {
    localStorage.removeItem(session_cookie_name);
}

document.addEventListener("DOMContentLoaded", function () {
    checkSessionActive();
});
