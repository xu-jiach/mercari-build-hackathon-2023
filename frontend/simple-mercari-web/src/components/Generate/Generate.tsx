import React, { useState } from 'react';
import { useCookies } from "react-cookie";

interface DescriptionGeneratorProps {
  onGenerated: (description: string) => void;
}

const DescriptionGenerator: React.FC<DescriptionGeneratorProps> = ({
  onGenerated,
}) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [generatedDescription, setGeneratedDescription] = useState("");
  const [cookies] = useCookies(["token", "userID"]);
  const [itemName, setItemName] = useState("");
  const [categoryID, setCategoryID] = useState(0);

  const handleGenerateDescription = async () => {
    try {
      setLoading(true);
      setError("");

      const response = await fetch("http://localhost:9000/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${cookies.token}`,
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
      setGeneratedDescription(responseData.description);
      onGenerated(responseData.description);
    } catch (error) {
      setError("An error occurred while generating the description.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Generate Description</h2>
      <div>
        <label htmlFor="itemName">Item Name:</label>
        <input
          type="text"
          id="itemName"
          value={itemName}
          onChange={(e) => setItemName(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="categoryID">Category ID:</label>
        <input
          type="number"
          id="categoryID"
          value={categoryID}
          onChange={(e) => setCategoryID(parseInt(e.target.value))}
        />
      </div>
      {loading && <p>Loading...</p>}
      {error && <p>{error}</p>}
      {generatedDescription && <p>{generatedDescription}</p>}
      <button onClick={handleGenerateDescription} disabled={loading}>
        Generate
      </button>
    </div>
  );
};

export default DescriptionGenerator;

