import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';

export const Login = () => {
  const [userID, setUserID] = useState<number>();
  const [password, setPassword] = useState<string>();
  const [_, setCookie] = useCookies(["userID", "token"]);

  const navigate = useNavigate();

  const onSubmit = (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    fetcher<{ id: number; name: string; token: string }>(`/login`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        user_id: userID,
        password: password,
      }),
    })

      .then((user) => {
        toast.success("Signed in!");
        console.log("POST success:", user.id);
        setCookie("userID", user.id);
        setCookie("token", user.token);
        navigate("/");
      })
      .catch((err) => {
        console.log(`POST error:`, err.status);
        toast.error("Error: " + err.status + " User id or Password incorrect");
      });

  };

  return (
    <div>
      <h2>Welcome back</h2>
      <div className="Login">
        <TextField id="outlined-basic" label="user id" variant="outlined" className="text-boxes"
          type="number"
          name="userID"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setUserID(Number(e.target.value));
          }}
          required
        /> <br/>
        <TextField id="outlined-basic" label="password" variant="outlined" className="text-boxes" sx={{ mt: 3}}
          type="password"
          name="password"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setPassword(e.target.value);
          }}
        />
        <Button variant="contained" onClick={onSubmit} id="sign-in-up-btn" color="error" sx={{ mt: 3}}>
          Login
        </Button>
      </div>
    </div>
  );
};
