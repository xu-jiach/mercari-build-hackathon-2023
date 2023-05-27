import React, { useEffect, useState } from "react";
import {useParams} from "react-router-dom";
import { useCookies } from "react-cookie";
import { MerComponent } from "../MerComponent";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import Button from '@mui/material/Button';

interface Category {
  id: number;
  name: string;
}

type formDataType = {
  name: string;
  category_id: number;
  newCategory: string;
  price: number;
  description: string;
  image: string | File;
};

export const Listing: React.FC = () => {
  const initialState = {
    name: "",
    category_id: 0,
    newCategory: "",
    price: 0,
    description: "",
    image: "",
  };
  const [values, setValues] = useState<formDataType>(initialState);
  const [categories, setCategories] = useState<Category[]>([]);
  // Add new state to handle new category
  const [newCategory, setNewCategory] = useState<string>("");
  // Add new state to handle new category checkbox
  const [newCategoryCheckboxChecked, setNewCategoryCheckboxChecked] = useState(false);
  const [cookies] = useCookies(["token", "userID"]);
  // Add an itemId state variable, it's null when creating a new item, set to the item's id when editing an existing item.
  const { itemId } = useParams<{ itemId: string }>();
  const isEditing = itemId !== undefined;


  const onValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value,
    });
  };

  const onSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value,
    });
  };

  const onFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.files![0],
    });
  };

  // This function will handle changes in the newCategory input
  const onNewCategoryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNewCategory(event.target.value);
  };

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData();
    data.append("name", values.name);
    data.append("category_id", values.category_id.toString());
    data.append("price", values.price.toString());
    data.append("description", values.description);
    data.append("image", values.image);

    // If category_id is 0, add newCategory to the FormData
    if (values.category_id === 0) {
      data.append("category_name", newCategory);
    }

    if (isEditing) {
      // Send a PUT request to update the existing item
      fetcher(`/items/${itemId}`, {
        method: "PUT",
        body: data,
        headers: {
          Authorization: `Bearer ${cookies.token}`,
        },
      })
        .then(() => {
          toast.success("Item updated successfully!");
        })
        .catch((error: Error) => {
          toast.error(error.message);
          console.error("PUT error:", error);
        });
    } else {
      // Send a POST request to create a new item
      fetcher(`/items`, {
        method: "POST",
        body: data,
        headers: {
          Authorization: `Bearer ${cookies.token}`,
        },
      })
        .then((res) => {
          sell(res.id);
        })
        .catch((error: Error) => {
          toast.error(error.message);
          console.error("POST error:", error);
        });
    }
  };

  const fetchItemDetails = () => {
    if (isEditing) {
      fetcher(`/items/${itemId}`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${cookies.token}`,
        },
      })
        .then((item) => {
          const matchingCategory = categories.find(
            (category) => category.id === item.category_id
          );

          setValues((prevValues) => ({
            ...prevValues,
            name: item.name,
            category_id: item.category_id,
            price: item.price,
            description: item.description,
            image: item.image, // assuming item.image is the URL of the image
          }));
        })
        .catch((error: Error) => {
          toast.error(error.message);
          console.error("GET error:", error);
        });
    }
  };



  const sell = (itemID: number) =>
    fetcher(`/sell`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
      body: JSON.stringify({
        item_id: itemID,
      }),
    })
      .then((_) => {
        toast.success("Item added successfully!");
      })
      .catch((error: Error) => {
        toast.error(error.message);
        console.error("POST error:", error);
      });

      const fetchCategories = () => {
        fetcher<Category[]>(`/items/categories`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
          },
        })
          .then((items) => setCategories(items))
          .catch((err) => {
            console.log(`GET error:`, err);
            toast.error(err.message);
          });
      };

      useEffect(() => {
        fetchCategories();
      }, []);


  // Effect that runs whenever the new category name changes
  useEffect(() => {
    const matchingCategory = categories.find(
      // avoid discrepanies between lowercase and uppercase
      (category) => category.name.toLowerCase() === newCategory.toLowerCase()
    );

    if (matchingCategory) {
      setValues({
        ...values,
        category_id: matchingCategory.id,
        newCategory: "",
      });
      setNewCategoryCheckboxChecked(false);
      // clear the new category input
      setNewCategory("");
    }
  }, [newCategory, categories]);




  return (
  <MerComponent>
    <div className="Listing">
      <form onSubmit={onSubmit} className="ListingForm">
        <div>
            <input
            type="text"
            name="name"
            id="MerTextInput"
            placeholder="name"
            onChange={onValueChange}
            value={values.name}
            required
          />
          <select
            name="category_id"
            id="MerTextInput"
            value={values.category_id}
            onChange={onSelectChange}
          >
            <option value="0">-</option>
            {categories &&
              categories.map((category) => {
                return <option value={category.id}>{category.name}</option>;
              })}
          </select>
            <input
            type="checkbox"
            id="newCategoryCheckbox"
            name="newCategoryCheckbox"
            checked={newCategoryCheckboxChecked}
            onChange={(event) => {
              const checked = event.target.checked;
              setNewCategoryCheckboxChecked(checked);
              if (checked) {
                setValues({
                  ...values,
                  category_id: 0,
                  newCategory: "",
                });
              } else {
                setValues({
                  ...values,
                  category_id: 0,
                  newCategory: "",
                });
              }
            }}
          />
          <label htmlFor="newCategoryCheckbox">Create a new category</label>
          <input
            type="text"
            name="newCategory"
            id="newCategory"
            placeholder="Enter new category"
            onChange={onNewCategoryChange}
            disabled={values.category_id !== 0}
        />
          <input
            type="number"
            name="price"
            id="MerTextInput"
            placeholder="price"
            onChange={onValueChange}
            value={values.price}
            required
          />
          <input
            type="text"
            name="description"
            id="MerTextInput"
            placeholder="description"
            onChange={onValueChange}
            value={values.description}
            required
          />
          <input
            type="file"
            name="image"
            id="MerTextInput"
            onChange={onFileChange}
            required
          />
          <Button variant="contained" type="submit" color="primary" sx={{ mt: 3}}>
            List
          </Button>
        </div>
      </form>
    </div>
  </MerComponent>

  );
};
