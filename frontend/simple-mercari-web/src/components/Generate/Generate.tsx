import React, { useState, useEffect } from 'react';

interface DescriptionGeneratorProps {
  itemName: string;
  categoryID: number;
  token: any;
  onGenerated: (description: string) => void;
}

const DescriptionGenerator: React.FC<DescriptionGeneratorProps> = ({
  itemName,
  categoryID,
  token,
  onGenerated,
}) => {
  useEffect(() => {
    generateDescription();
  }, [itemName, categoryID, token]);

  const generateDescription = async () => {
    try {
      const response = await fetch("http://localhost:9000/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          itemName,
          categoryID,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to generate description");
      }

      const responseData = await response.json();
      const description = responseData.description;
      onGenerated(description);
    } catch (error) {
      console.error("An error occurred while generating the description:", error);
    }
  };

  return null; // Since we don't render any elements in this component
};

export default DescriptionGenerator;
