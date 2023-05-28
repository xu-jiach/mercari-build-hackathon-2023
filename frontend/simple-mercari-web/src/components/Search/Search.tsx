import {MerComponent} from "../MerComponent";
import {ItemList} from "../ItemList";
import React from "react";
import {Item} from "../../common/interfaces";
import {fetcher} from "../../helper";
import {useEffect, useState} from "react";
import "./search.css"
import { Typography } from '@mui/material';
import theme from "../../theme";

export const Search = () => {
    const [items, setItems] = useState<Item[]>([]);
    const [keyword, setKeyword] = useState<string>("");
    const query = new URLSearchParams(window.location.search);

    console.log("keyword", keyword)
    const fetchItems = () => {
        fetcher<Item[]>(`/search?name=` + keyword , {
            method: "GET",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
        })
            .then((res) => {
                console.log("GET success:", res);
                setItems(res);
            })

            .catch((err) => {
                console.log(`GET error:`, err);
            });
    }

    useEffect(() => {
        console.log("fetching items");
        fetchItems();
    }, [keyword]);

    useEffect(() => {
        setKeyword(query.get("keyword") || "");
    }, []);

    return (
        <MerComponent>
            <div className={"search-container"}>
                <div className={"search-title"}>
                    <Typography variant="h4" component="p" color={theme.palette.common.black}>
                        {/* TODO: Display some message when the query is empty */}
                        {keyword ? `Search results for "${keyword}"` : ""}
                        <span> </span>
                        <Typography variant="h6" component="span" color={theme.palette.grey["600"]}>
                            ({items?.length ?? 0} results)
                        </Typography>
                    </Typography>
                </div>
                {/* TODO: Distinguish items that are sold out and still for sale */}
                <ItemList items={items} />
            </div>
        </MerComponent>
    )
}