import { useParams } from 'react-router-dom';
import { useEffect, useState } from "react";
import { fetcher } from "../../helper";
import { toast } from "react-toastify";
import { ItemList } from "../ItemList";
import { useCookies } from "react-cookie";
import { Categories } from '../Categories/Categories';
import { Footer } from '../Footer';

interface Item {
  id: number;
  name: string;
  price: number;
  category_name: string;
}

export const CategoryPage: React.FC = () => {
  const [items, setItems] = useState<Item[]>([]);
  const { id } = useParams();
  const [cookies] = useCookies(["token"]);

  const fetchItems = () => {
    fetcher<Item[]>(`/categories/${id}/items`, {
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

  useEffect(() => {
    fetchItems();
  }, [id]);

  return (
    <div>
      <Footer/>
      <Categories/>
      <ItemList items={items} />
    </div>
  );
};
