
export const ItemImage: React.FC<{ itemImage: Blob }>  = ({itemImage}) => {
  return (
      <div className={"image-container"}>
          <img
              height={480}
              width={480}
              src={URL.createObjectURL(itemImage)}
              alt="thumbnail"
              className={"item-image"}
          />
      </div>
  );
}