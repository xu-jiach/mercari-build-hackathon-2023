import React from "react";
import { Item } from "../Item";
import "./ItemList.css";
import { Item as ItemInterface } from "../../common/interfaces";

interface Prop {
  items: ItemInterface[];
}

export const ItemList: React.FC<Prop> = (props) => {
  return (
    <section className={"item-container"}>
      <div className={"item-list"}>
        {props.items &&
          props.items.map((item) => {
            return <Item item={item} />;
          })}
      </div>
    </section>
  );
};
