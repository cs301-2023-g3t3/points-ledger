import http from "k6/http";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "1s",
  ext: {
    loadimpact: {
      // Project: Default project
      projectID: 3667423,
      // Test runs with the same name groups test runs together.
      name: "Load Test",
    },
  },
};

export default function () {
  const params = {
    headers: {
      Authorization: `Bearer ${__ENV.AUTH_TOKEN}`,
    },
  };
  http.get(
    "https://vm1swtn9ii.execute-api.ap-southeast-1.amazonaws.com/v1/users/roles",
    params
  );
  sleep(1);
}
