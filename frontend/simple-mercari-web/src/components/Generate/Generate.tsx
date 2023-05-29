import React, { useState } from 'react';
import { useCookies } from "react-cookie";

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
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [generatedDescription, setGeneratedDescription] = useState("");
  const [cookies] = useCookies(["token", "userID"]);
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
      {/* Render loading and error messages */}
      {loading && <p>Loading...</p>}
      {error && <p>{error}</p>}
      {/* Render the generated description */}
      {generatedDescription && <p>{generatedDescription}</p>}
      {/* Button to trigger description generation */}
      <button onClick={handleGenerateDescription} disabled={loading}>
        Generate
      </button>
    </div>
  );
};

export default DescriptionGenerator;
