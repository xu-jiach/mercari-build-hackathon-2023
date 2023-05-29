import React, { useState } from "react";
import "./ItemImage.css";

export const ItemImage: React.FC<{ itemImage: Blob }> = ({ itemImage }) => {
  const [isZoomed, setIsZoomed] = useState(false);
  const [cursorPosition, setCursorPosition] = useState({ x: 0, y: 0 });

  const handleMouseMove = (event: React.MouseEvent<HTMLDivElement>) => {
    const { left, top, width, height } = event.currentTarget.getBoundingClientRect();
    const x = ((event.clientX - left) / width) * 100;
    const y = ((event.clientY - top) / height) * 100;
    setCursorPosition({ x, y });
  };

  const handleMouseEnter = () => {
    setIsZoomed(true);
  };

  const handleMouseLeave = () => {
    setIsZoomed(false);
  };

  const imageContainerClasses = `image-container ${isZoomed ? "zoomed" : ""}`;
  const imageClasses = `item-image ${isZoomed ? "zoomed" : ""}`;
  const zoomStyles = {
    backgroundImage: `url(${URL.createObjectURL(itemImage)})`,
    backgroundPosition: `${cursorPosition.x}% ${cursorPosition.y}%`,
    backgroundSize: '200%', // double the size for a zoom effect
    backgroundRepeat: 'no-repeat',
  };

  return (
    <div
      className={imageContainerClasses}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onMouseMove={handleMouseMove}
    >
      <div className="zoomed-image" style={isZoomed ? zoomStyles : {}} />
      <img src={URL.createObjectURL(itemImage)} alt="thumbnail" className={imageClasses} />
    </div>
  );
};
