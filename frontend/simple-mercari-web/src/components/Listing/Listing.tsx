import React, { useEffect, useState, ReactNode, ChangeEvent}  from "react";
import {useParams, useNavigate} from "react-router-dom";
import { useCookies } from "react-cookie";
import { MerComponent } from "../MerComponent";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import { Button, TextField, FormControl, InputLabel, Select, MenuItem, Checkbox, FormControlLabel, SelectChangeEvent } from '@mui/material';
import Tooltip, { TooltipProps, tooltipClasses } from '@mui/material/Tooltip';


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
  item_passcode: string;
};

export const Listing: React.FC = () => {
  const initialState = {
    name: "",
    category_id: 0,
    newCategory: "",
    price: 0,
    description: "",
    image: "",
    item_passcode: "",
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
  const [fileName, setFileName] = useState("");
  const isEditing = itemId !== undefined;
  const navigate = useNavigate();
  const inPersonDescription = "Passcode required, only share with in person buyer"


  const onValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value,
    });
  };

  const onSelectChange = (event: SelectChangeEvent<number>) => {
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
    setFileName(event.target.files![0]?.name || ""); // Set file name
  };

  // This function will handle changes in the newCategory input
  const onNewCategoryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNewCategory(event.target.value);
  };

  // call the addcategory function when the checkbox is checked
  const addCategory = async () => {
    try {
      const response = await fetcher(`/categories`, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
          Authorization: `Bearer ${cookies.token}`,
        },
        body: JSON.stringify({
          name: newCategory,
        }),
      });
      toast.success("New category created successfully!");
      return response.id;
    } catch (error: any) {
      toast.error(error.message);
      console.error("POST error:", error);
      return null;
    }
  };


  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData();
    data.append("name", values.name);
    data.append("price", values.price.toString());
    data.append("description", values.description);
    data.append("image", values.image);
    data.append("item_password", values.item_passcode);

  // If category_id is 0, create a new category and get its id
    if (values.category_id === 0) {
      const categoryId = await addCategory();
    if (!categoryId) {
      toast.error("Failed to create new category. Please try again.");
      return;
    }
    data.append("category_id", categoryId.toString());
  } else {
    data.append("category_id", values.category_id.toString());
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
          sell(Number(itemId), isEditing);
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
          sell(res.id, isEditing);
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

    const sell = (itemID: number, isEditing: boolean) =>
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
        // only display "Item added successfully!" toast if not editing
        if (!isEditing) {
          toast.success("Item added successfully!");
        }
        navigate('/'); // Redirect to homepage
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
        fetchItemDetails();
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

  const [allowInPersonPurchases, setAllowInPersonPurchases] = useState(false);



  return(
    <MerComponent>
      <div className="Listing">
        <h1>List a new item</h1>
        <form onSubmit={onSubmit}>
          <div className="listing-form">
              <TextField
                id="name"
                name="name"
                value={values.name}
                onChange={onValueChange}
                label="Name"
                required
                sx={{ mt: 3, mb: 3}}
              />
              <FormControl>
                <InputLabel id="category-label">Category</InputLabel>
                <Select
                  labelId="category-label"
                  id="category_id"
                  name="category_id"
                  value={values.category_id}
                  onChange={onSelectChange}
                  sx={{ mt: 3 }}
                >
                  <MenuItem value={0} disabled>Select a category</MenuItem>
                  {categories &&
                    categories.map((category) => {
                      return <MenuItem value={category.id}>{category.name}</MenuItem>;
                    })}
                </Select>
              </FormControl>
              <FormControlLabel sx={{ mt: 3 }}
                control={
                  <Checkbox
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
                }
                label="Create a new category"
              />
              <TextField
                id="newCategory"
                name="newCategory"
                value={newCategory}
                onChange={onNewCategoryChange}
                label="New Category"
                disabled={!newCategoryCheckboxChecked}
                sx={{ mt: 3 }}
              />
              <TextField
                type="number"
                id="price"
                name="price"
                value={values.price}
                onChange={onValueChange}
                label="Price"
                required
                sx={{ mt: 3 }}
              />
              <TextField
                id="description"
                name="description"
                value={values.description}
                onChange={onValueChange}
                label="Description"
                required
                multiline
                rows={4}
                sx={{ mt: 3 }}
              />
              <Button variant="contained" component="label" sx={{ mt: 3 }}>
                Upload Image
                <input
                  type="file"
                  name="image"
                  id="image"
                  onChange={onFileChange}
                  required
                  hidden
                />
              </Button>
              {fileName && <div className="mt1">Selected file: {fileName}</div>}

              <Tooltip title={inPersonDescription} arrow>
                <FormControlLabel sx={{ mt: 3 }}
                  control={
                    <Checkbox
                      checked={allowInPersonPurchases}
                      onChange={(event) => setAllowInPersonPurchases(event.target.checked)}
                    />
                  }
                  label="Allow in person purchases"
                />
              </Tooltip>
              <TextField
                id="item_passcode"
                name="item_passcode"
                value={values.item_passcode}
                onChange={onValueChange}
                label="In Person Passcode"
                sx={{ mt: 3 }}
                disabled={!allowInPersonPurchases}
              />
              <Button variant="contained" type="submit" color="secondary" sx={{ mt: 3 }}>
                List
              </Button>
          </div>
        </form>
      </div>
    </MerComponent>
  )

};
