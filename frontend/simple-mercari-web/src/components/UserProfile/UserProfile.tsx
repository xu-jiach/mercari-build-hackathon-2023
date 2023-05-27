import { useState, useEffect } from "react";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";
import { MerComponent } from "../MerComponent";
import { toast } from "react-toastify";
import { ItemList } from "../ItemList";
import { fetcher } from "../../helper";

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
  const [cookies] = useCookies(["token"]);
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

  return (
    <MerComponent>
      <div className="UserProfile">
        <div>
          <div>
            <h2>
              <span>Balance: {balance}</span>
            </h2>
            <input
              type="number"
              name="balance"
              id="MerTextInput"
              placeholder="0"
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setAddedBalance(Number(e.target.value));
              }}
              required
            />
            <button onClick={onBalanceSubmit} id="MerButton">
              Add balance
            </button>
          </div>

          <div>
            <h2>Item List</h2>
            {<ItemList items={items} />}
          </div>
        </div>
      </div>
    </MerComponent>
  );
};
