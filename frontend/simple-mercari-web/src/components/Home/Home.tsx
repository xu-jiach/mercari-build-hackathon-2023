import { Login } from "../Login";
import { Signup } from "../Signup";
import { ItemList } from "../ItemList";
import { useCookies } from "react-cookie";
import { MerComponent } from "../MerComponent";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import "react-toastify/dist/ReactToastify.css";
import Divider from '@mui/material/Divider';
import { styled } from '@mui/material/styles';
import {Item as ItemInterface} from "../../common/interfaces";

export const Home = () => {
  const [cookies] = useCookies(["userID", "token"]);
  const [items, setItems] = useState<ItemInterface[]>([]);

  const fetchItems = () => {
    fetcher<ItemInterface[]>(`/items`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((data) => {
        console.log("GET success:", data);
        setItems(data);
      })
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error("Error: " + err.status);
      });
  };

  useEffect(() => {
    fetchItems();
  }, []);

  const Root = styled('div')(({ theme }) => ({
    width: "300px",
    ...theme.typography.body2,
    '& > :not(style) + :not(style)': {
      marginTop: theme.spacing(2),
    },
  }));

  const signUpAndSignInPage = (
    <div className="sign-in-up-page">
      <div id="account-form">
        <div>
          <Login />
        </div>
        <Root><Divider sx={{ mt: 3}}>or</Divider>  </Root>
        <div>
          <Signup />
        </div>
      </div>
      <div id="featured-items">
        <img id="featured-items-img" src="https://media.wired.com/photos/629133e5e9a46d033b3380c7/master/w_2560%2Cc_limit/Finding-a-PlayStation-5-Is-About-to-Get-Easier-Gear-shutterstock_1855958302.jpg" alt="(placeholder)featured" />
      </div>
    </div>

  );

  const itemListPage = (
    <MerComponent>
      <div>
        <ItemList items={items} />
      </div>
    </MerComponent>
  );

  return <>{cookies.token ? itemListPage : signUpAndSignInPage}</>;
};