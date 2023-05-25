import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';

export const Signup = () => {
  const [name, setName] = useState<string>();
  const [password, setPassword] = useState<string>();
  const [userID, setUserID] = useState<number>();
  const [_, setCookie] = useCookies(["userID"]);

  const navigate = useNavigate();
  const onSubmit = (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    fetcher<{ id: number; name: string }>(`/register`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        password: password,
      }),
    })
      .then((user) => {
        toast.success("New account is created!");
        console.log("POST success:", user.id);
        setCookie("userID", user.id);
        setUserID(user.id);
        navigate("/");
      })
      .catch((err) => {
        console.log(`POST error:`, err);
        toast.error(err.message);
      });
  };

  return (
    <div>
      <div className="Signup">
        <TextField id="outlined-basic" label="username" variant="outlined" className="text-boxes" sx={{ mt: 3}}
          type="text"
          name="name"
          placeholder="name"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setName(e.target.value);
          }}
          required
        /> <br/>
        <TextField id="outlined-basic" label="password" variant="outlined" className="text-boxes" sx={{ mt: 3}}
          type="password"
          name="password"
          placeholder="password"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setPassword(e.target.value);
          }}
        />
        <Button variant="outlined" onClick={onSubmit} id="sign-in-up-btn" color="error" sx={{ mt: 3}}>
          Signup
        </Button>
        {userID ? (
          <p>Use "{userID}" as UserID for login</p>
        ) : null}
      </div>
    </div>
  );
};
