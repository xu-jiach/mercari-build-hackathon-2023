import { Item, ItemStatus } from "../../common/interfaces";
import { useNavigate, useParams } from "react-router-dom";
import React, { useEffect, useState, ReactNode, ChangeEvent} from "react";
import { useCookies } from "react-cookie";
import { fetcher } from "../../helper";
import { toast } from "react-toastify";
import Button from "@mui/material/Button";
import { Chip, TextField, FormControlLabel, Checkbox } from "@mui/material";

type InPersonPasscode = {
    password: string;
}

type IsInPersonAvailable = {
    isAvailable: boolean;
}

export const ItemDescription: React.FC<{ item: Item, isOwner: boolean}>  = ({item, isOwner}) => {
    const [imWithSeller, setImWithSeller] = useState(false);
    const [inPersonPasscode, setInPersonPasscode] = useState<string | null>(null);
    const [isInPersonAvailable, setIsInPersonAvailable] = useState<boolean>(false);

    const navigate = useNavigate();
    const [cookies] = useCookies(["token", "userID"]);
    const params = useParams();

    const onSubmit = (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        fetcher<Item[]>(`/purchase/${params.id}`, {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
                Authorization: `Bearer ${cookies.token}`,
            },
            body: JSON.stringify({
                user_id: Number(cookies.userID),
            }),
        })
            .then(() => window.location.reload())
            .catch((err) => {
                console.log(`POST error:`, err);
                toast.error("Error: " + err.message);
            });
    };

    const getInPersonPasscode = () => {
        fetcher<InPersonPasscode>(`/items/${params.id}/pass`, {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
                Authorization: `Bearer ${cookies.token}`,
            },
        })
            .then((res) => {
                console.log(`GET response:`, res);
                setInPersonPasscode(res.password);
            })
            .catch((err) => {
                console.log(`GET error:`, err);
            });
    }

    const getIsInPersonAvailable = () => {
        fetcher<IsInPersonAvailable>(`/onsite-purchase/${params.id}/available`, {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
                Authorization: `Bearer ${cookies.token}`,
            },
        })
            .then((res) => {
                setIsInPersonAvailable(res.isAvailable);
            })
            .catch((err) => {
                console.log(`GET error:`, err);
            });
    }



    useEffect(() => {
        if (isOwner) {
            getInPersonPasscode();
        } else {
            getIsInPersonAvailable();
        }
    }, []);

    return (
        <div className={"description-container"}>
            <h1>{item.name}</h1>
            <h2><span id={"yen-symbol"}>¥</span> {item.price.toLocaleString()}</h2>
            <Chip label={item.category_name} component="a" /> {/* TODO: Navigate to category view on clicking */}

            {item.status == ItemStatus.ItemStatusSoldOut ? (
                <Button disabled={true} onClick={onSubmit} id="MerDisableButton">
                    SoldOut
                </Button>
            ) : (
                <>
                    {isOwner && (
                      <>
                        <Button
                            onClick={() => navigate(`/edit-item/${item.id}`)} // Navigate to /edit-item/:itemId when the Edit button is clicked
                            id="MerButton"
                        >
                            Edit
                        </Button>
                         {inPersonPasscode &&(
                          <p><strong>In person purchasing with passcode: {inPersonPasscode}</strong></p>
                        )}
                      </>
                    )}
                    <hr/>
                    {!isOwner && (
                      <>
                        <Button variant="contained" onClick={onSubmit} id="buy-now-btn" color="primary" sx={{ mt: 3}}>
                          Buy now
                        </Button>
                          {isInPersonAvailable && (
                              <>
                                  <FormControlLabel sx={{mt: 3}}
                                                  control={<Checkbox
                                                      checked={imWithSeller}
                                                      onChange={(event) => setImWithSeller(event.target.checked)}/>}
                                                  label="I'm with the owner"/><TextField sx={{mt: 2, ml: 3}}
                                                                                         id="inPersonKey"
                                                                                         name="inPersonKey"
                                                                                         label="In Person Passcode"
                                                                                         disabled={!imWithSeller}/>
                              </>
                          )}
                      </>
                    )}




                </>
            )}

            <p>{item.description}</p>
            <p>User: {item.user_id}</p> {/* TODO: Display user name instead of user id */}
        </div>
    );
}
