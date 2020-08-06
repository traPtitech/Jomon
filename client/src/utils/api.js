import axios from "axios";

axios.defaults.withCredentials = true;
export const traQBaseURL = "https://q.trap.jp/api/v3";
// axios.defaults.baseURL = "http://localhost:8080";
//   process.env.NODE_ENV === "development"
//     ? "http://localhost:3000"
//     : process.env.VUE_APP_API_ENDPOINT;

export async function redirectAuthEndpoint() {
  try {
    const data = (await axios.get("/api/auth/genpkce")).data;

    const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`);

    authorizationEndpointUrl.search = new URLSearchParams({
      responseType: "code",
      clientId: data.clientId,
      state: data.state,
      codeChallenge: data.codeChallenge,
      codeChallengeMethod: data.codeChallengeMethod
    }).toString();

    window.location.assign(authorizationEndpointUrl.toString());
  } catch (err) {
    console.log(err);
  }
}
