import { useState, useEffect } from "react";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router-dom";
import { fetcherBlob } from "../../helper";
import "./Item.css";
import { Item as ItemInterface } from "../../common/interfaces";

export const Item: React.FC<{ item: ItemInterface }> = ({ item }) => {
  const navigate = useNavigate();
  const [itemImage, setItemImage] = useState<string>("");
  const [cookies] = useCookies(["token"]);

  async function getItemImage(itemId: number): Promise<Blob> {
    return await fetcherBlob(`/items/${itemId}/image`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
    });
  }

  useEffect(() => {
    async function fetchData() {
      const image = await getItemImage(item.id);
      setItemImage(URL.createObjectURL(image));
    }

    fetchData();
  }, [item]);

  return (
      <div className={`item ${item.status == 3 ? "sold" : ""}`}>
        <img
          src={itemImage}
          alt={item.name}
          height={180}
          width={180}
          onClick={() => navigate(`/item/${item.id}`)}
        />
        <div className={"price-container"}>
          <span className={"price"}>¥{item.price.toLocaleString()}</span>
        </div>
      </div>
  );
};
