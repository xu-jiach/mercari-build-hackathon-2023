import { useState, useEffect } from "react";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";
import { MerComponent } from "../MerComponent";
import { toast } from "react-toastify";
import { ItemList } from "../ItemList";
import { fetcher } from "../../helper";
import { Button, InputAdornment, TextField, FormControl, InputLabel, Select, MenuItem, Checkbox, FormControlLabel, SelectChangeEvent } from '@mui/material';
import { FaPlusCircle} from 'react-icons/fa';
import ListItemText from '@mui/material/ListItemText';
import OutlinedInput from '@mui/material/OutlinedInput';

interface Item {
  id: number;
  name: string;
  price: number;
  category_name: string;
}

export const UserProfile: React.FC = () => {
  const [items, setItems] = useState<Item[]>([]);
  const [balance, setBalance] = useState<number>();
  const [addedbalance, setAddedBalance] = useState<number>();
  const [cookies] = useCookies(["userID", "token"])
  const params = useParams();

  const fetchItems = () => {
    fetcher<Item[]>(`/users/${params.id}/items`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
    })
      .then((items) => setItems(items))
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error("Error: " + err.status);
      });
  };

  const fetchUserBalance = () => {
    fetcher<{ balance: number }>(`/balance`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
    })
      .then((res) => {
        setBalance(res.balance);
      })
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error("Error: " + err.status);
      });
  };

  useEffect(() => {
    fetchItems();
    fetchUserBalance();
  }, []);

  const onBalanceSubmit = () => {
    fetcher(`/balance`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
      body: JSON.stringify({
        balance: addedbalance,
      }),
    })
      .then((_) => window.location.reload())
      .catch((err) => {
        console.log(`POST error:`, err);
        toast.error("Error: " + err.status);
      });
  };

  const formattedBalance = balance?.toLocaleString(); // Format balance with commas

  return (
    <MerComponent>
      <div className="UserProfile">
        <h2>Wallet</h2>
        <div>
          <div className="wallet-balance">
            <p><strong>Total Balance: </strong></p>
            <div className="balance-display">
            <h2>¥{formattedBalance}</h2>
            <FormControl
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setAddedBalance(Number(e.target.value));
              }}>
              <InputLabel htmlFor="outlined-adornment-amount">Amount</InputLabel>
              <OutlinedInput
                id="outlined-adornment-amount"
                startAdornment={<InputAdornment position="start">¥</InputAdornment>}
                label="Amount"
              />
            </FormControl>
            <button onClick={onBalanceSubmit}>
              <FaPlusCircle color="#FF5757" />
              <ListItemText primary="Add Balance" />
            </button>
            </div>
          </div>

          <div className="transactions-and-listings">
            <div className="transactions">
              <h2>Recent Transactions</h2>
              <div className="transaction-unit">
                <p>Deposit </p>
                <p className="deposit">+800¥</p>
              </div>
              <div className="transaction-unit">
              <p>Purchase </p>
                <p className="purchase">-8000¥</p>
              </div>
            </div>
            <div className="listings">
              <h2>Your listings</h2>
              {<ItemList items={items} />}
            </div>
          </div>

        </div>
      </div>
    </MerComponent>
  );
};
