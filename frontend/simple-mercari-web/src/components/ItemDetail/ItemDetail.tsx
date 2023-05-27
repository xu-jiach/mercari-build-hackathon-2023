import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useCookies } from "react-cookie";
import { MerComponent } from "../MerComponent";
import { toast } from "react-toastify";
import { fetcher, fetcherBlob } from "../../helper";
import { ItemImage } from "./ItemImage";
import "./itemDetail.css";
import { Item } from "../../common/interfaces";
import { ItemDescription } from "./ItemDescription";

export const ItemDetail = () => {
    const [item, setItem] = useState<Item>();
    const [itemImage, setItemImage] = useState<Blob>();
    const [cookies] = useCookies(["token", "userID"]);
    const [isOwner, setIsOwner] = useState(false);
    const params = useParams();

  const fetchItem = () => {
    fetcher<Item>(`/items/${params.id}`, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    })
        .then((res) => {
          console.log("GET success:", res);
          setItem(res);
          setIsOwner(res.user_id === Number(cookies.userID)); // Set isOwner state here
        })

        .catch((err) => {
          console.log(`GET error:`, err);
          toast.error("Error: " + err.status);
        });

    fetcherBlob(`/items/${params.id}/image`, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    })
        .then((res) => {
          console.log("GET success:", res);
          setItemImage(res);
        })
        .catch((err) => {
          console.log(`GET error:`, err);
          toast.error("Error: " + err.status);
        });
  };

  useEffect(() => {
    fetchItem();
  }, []);

  return (
      <MerComponent condition={() => item !== undefined}>
        {item && itemImage && (
            <div className="item-detail-container">
                <ItemImage itemImage={itemImage} />
                <ItemDescription item={item} isOwner={isOwner}/>
            </div>
        )}
      </MerComponent>
  );
};
