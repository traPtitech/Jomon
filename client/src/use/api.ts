import axios from "axios";

axios.defaults.withCredentials = true;
export const traQBaseURL = "https://q.trap.jp/api/v3";
// axios.defaults.baseURL = "http://localhost:8080";
//   process.env.NODE_ENV === "development"
//     ? "http://localhost:3000"
//     : process.env.VUE_APP_API_ENDPOINT;

export async function redirectAuthEndpoint(): Promise<void> {
  const data = (await axios.get("/api/auth/genpkce")).data;
  const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`);

  authorizationEndpointUrl.search = new URLSearchParams({
    response_type: "code",
    client_id: data.client_id,
    code_challenge: data.code_challenge,
    code_challenge_method: data.code_challenge_method
  }).toString();

  window.location.assign(authorizationEndpointUrl.toString());
}
