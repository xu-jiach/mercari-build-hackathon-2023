import { server } from "./common/constants";

const wrap = <T>(task: Promise<Response>): Promise<T> => {
  return new Promise((resolve, reject) => {
    task
      .then((response) => {
        if (response.ok) {
          response
            .json()
            .then((json) => {
              resolve(json);
            })
            .catch((error) => {
              reject(error);
            });
        } else {
            response
                .json()
                // Return message in json to get the error message from the server
                .then((json) => {
                    const error = new Error();
                    error.name = `${response.status}: ${response.statusText}`;
                    error.message = json.message;
                    reject(error);
                })
                .catch((error) => {
                    reject(error);
                });
        }
      })
      .catch((error) => {
        reject(error);
      });
  });
};

const wrapBlob = (task: Promise<Response>): Promise<Blob> => {
  return new Promise((resolve, reject) => {
    task
      .then((response) => {
        if (response.ok) {
          response
            .blob()
            .then((blob) => {
              // jsonが取得できた場合だけresolve
              resolve(blob);
            })
            .catch((error) => {
              reject(error);
            });
        } else {
          reject(new Error(`${response.status}: ${response.statusText}`));
        }
      })
      .catch((error) => {
        reject(error);
      });
  });
};

export const fetcher = <T = any>(
  url: string,
  init?: RequestInit
): Promise<T> => {
  return wrap<T>(fetch(server.concat(url), init));
};

export const fetcherBlob = (url: string, init?: RequestInit): Promise<Blob> => {
  return wrapBlob(fetch(server.concat(url), init));
};
