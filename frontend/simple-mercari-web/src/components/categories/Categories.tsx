import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import { FaHamburger, FaTshirt, FaBed } from 'react-icons/fa';
import ListItemText from '@mui/material/ListItemText';


export const Categories = () => {
  const navigate = useNavigate();

  return(
    <div className="categories-container">
      <button className="category-icon" onClick={() => navigate('/categories/1')}>
        <FaHamburger color="#FF5757" />
        <ListItemText primary="Food" />
      </button>
      <button className="category-icon" onClick={() => navigate('/categories/2')}>
        <FaTshirt color="#FF5757" />
        <ListItemText primary="Fashion" />
      </button>
      <button className="category-icon" onClick={() => navigate('/categories/3')}>
        <FaBed color="#FF5757" />
        <ListItemText primary="Furniture" />
      </button>
    </div>
  )
}
