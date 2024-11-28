// REFERENCES:
// - PKCE Generator taken from: https://example-app.com/pkce

const CLIENT_ID = "beetle-quest";

function generateRandomString() {
    var array = new Uint32Array(56 / 2);
    window.crypto.getRandomValues(array);
    return Array.from(array, dec2hex).join("");
}

function dec2hex(dec) {
    return ("0" + dec.toString(16)).substr(-2);
}

function sha256(plain) {
    const encoder = new TextEncoder();
    const data = encoder.encode(plain);
    return window.crypto.subtle.digest("SHA-256", data);
}

function base64_urlencode(str) {
    return btoa(String.fromCharCode.apply(null, new Uint8Array(str)))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "");
}

async function pkce_challenge_from_verifier(v) {
    hashed = await sha256(v);
    var base64encoded = base64_urlencode(hashed);
    return base64encoded;
}

async function submitAuthorizationRequest(event) {
    event.preventDefault();

    let codeVerifier = generateRandomString();
    localStorage.setItem("codeVerifier", codeVerifier);

    let state = generateRandomString();
    localStorage.setItem("state", state);

    codeChallenge = await pkce_challenge_from_verifier(codeVerifier);

    const selectedScopes = Array.from(document.querySelectorAll('#authorizationForm input[name="scope"]:checked'))
        .map((input) => input.value)
        .join(", ");

    const reqBody = new FormData();
    reqBody.append("response_type", "code");
    reqBody.append("client_id", CLIENT_ID);
    reqBody.append("redirect_uri", "/api/v1/auth/tokenPage");
    reqBody.append("scope", selectedScopes);
    reqBody.append("state", state);
    reqBody.append("code_challenge", codeChallenge);
    reqBody.append("code_challenge_method", "S256");

    const xhr = event.detail.xhr;
    xhr.send(reqBody);
}

async function submitTokenRequest(event) {
    event.preventDefault();

    let state = localStorage.getItem("state");
    let codeVerifier = localStorage.getItem("codeVerifier");

    const code = document.getElementById("code").textContent.trim();
    const recvState = document.getElementById("state").textContent.trim();

    if (state !== recvState) {
        console.log("Invalid state!");
        return;
    }

    const reqBody = new FormData();
    reqBody.append("grant_type", "authorization_code");
    reqBody.append("code", code);
    reqBody.append("redirect_uri", "/api/v1/auth/tokenPage");
    reqBody.append("client_id", CLIENT_ID);
    reqBody.append("code_verifier", codeVerifier);

    const xhr = event.detail.xhr;
    xhr.send(reqBody);

    localStorage.removeItem("state");
    localStorage.removeItem("codeVerifier");
}

async function processTokenRequestResponse(event) {
    if (event.detail.xhr.status === 200) {
        const json_data = JSON.parse(event.detail.xhr.response);
        localStorage.setItem("ACCESS_TOKEN", json_data.access_token);

        window.location.href = "/static/";
    }
}

window.submitAuthorizationRequest = submitAuthorizationRequest;
window.submitTokenRequest = submitTokenRequest;
window.processTokenRequestResponse = processTokenRequestResponse;

function addAuthorzationHeader(event) {
    console.log(event);
    let tok = localStorage.getItem("ACCESS_TOKEN");
    if (tok) {
        event.detail.headers = {
            Authorization: `Bearer ${tok}`,
        };
    }
}

// document.addEventListener("htmx:config", (event) => {
//     const originalHeaders = htmx.config.headers || {};
//     htmx.config.headers = {
//         ...originalHeaders,
//         Authorization: `Bearer ${localStorage.getItem("ACCESS_TOKEN") || ""}`,
//     };
// });

document.addEventListener("htmx:configRequest", addAuthorzationHeader);
