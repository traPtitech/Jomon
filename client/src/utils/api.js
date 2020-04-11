// import axios from "axios";
import { randomString } from "./hash";

// axios.defaults.withCredentials = true
export const traQBaseURL = "https://q.trap.jp/api/1.0";
// axios.defaults.baseURL = "http://localhost:8080";
//   process.env.NODE_ENV === "development"
//     ? "http://localhost:3000"
//     : process.env.VUE_APP_API_ENDPOINT;

export async function redirectAuthorizationEndpoint(
  client_id,
  response_type,
  code_challenge,
  code_challenge_method
) {
  const state = randomString(10);
  sessionStorage.setItem(`state`, state);
  const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`);
  authorizationEndpointUrl.search = new URLSearchParams({
    client_id: client_id,
    response_type: response_type,
    code_challenge: code_challenge,
    code_challenge_method: code_challenge_method,
    state: state
  });
  window.location.assign(authorizationEndpointUrl);
}
