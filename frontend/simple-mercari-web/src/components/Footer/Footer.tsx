import * as React from 'react';
import { FaHome, FaCamera, FaUser } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import IconButton from '@mui/material/IconButton';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Toolbar from '@mui/material/Toolbar';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import MenuIcon from '@mui/icons-material/Menu';

export const Footer: React.FC = () => {
  const [cookies] = useCookies(['userID']);
  const navigate = useNavigate();
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  if (!cookies.userID) {
    return <></>;
  }

  return (
    <>
      <footer>
        <div className="MerFooterItem" onClick={() => navigate('/')}>
          <FaHome />
          <p>Home</p>
        </div>
        <div className="MerFooterItem" onClick={() => navigate('/sell')}>
          <FaCamera />
          <p>Listing</p>
        </div>
        <div
          className="MerFooterItem"
          onClick={() => navigate(`/user/${cookies.userID}`)}
        >
          <FaUser />
          <p>MyPage</p>
        </div>
      </footer>
      <Drawer variant="permanent" open={open}>
        <Toolbar>
          <IconButton
            onClick={handleDrawerClose}
            sx={{
              ...(open && { display: 'none' }),
            }}
          >
            <ChevronLeftIcon />
          </IconButton>
        </Toolbar>
        <List>
          <ListItem button onClick={() => navigate('/')}>
            <ListItemIcon>
              <FaHome />
            </ListItemIcon>
            <ListItemText primary="Home" />
          </ListItem>
          <ListItem button onClick={() => navigate('/sell')}>
            <ListItemIcon>
              <FaCamera />
            </ListItemIcon>
            <ListItemText primary="Listing" />
          </ListItem>
          <ListItem
            button
            onClick={() => navigate(`/user/${cookies.userID}`)}
          >
            <ListItemIcon>
              <FaUser />
            </ListItemIcon>
            <ListItemText primary="MyPage" />
          </ListItem>
        </List>
      </Drawer>
      <Box component="div" sx={{ flexGrow: 1, p: 3 }}>
        <IconButton
          color="inherit"
          aria-label="open drawer"
          onClick={handleDrawerOpen}
          edge="start"
          sx={{
            ...(open && { display: 'none' }),
          }}
        >
          <MenuIcon />
        </IconButton>
        <Typography variant="h6" noWrap component="div">
          Mini variant drawer
        </Typography>
      </Box>
    </>
  );
};

// import { FaHome, FaCamera, FaUser } from "react-icons/fa";
// import { useNavigate } from "react-router-dom";
// import { useCookies } from "react-cookie";
// import "./Footer.css";

// export const Footer: React.FC = () => {
//   const [cookies] = useCookies(["userID"]);
//   const navigate = useNavigate();

//   if (!cookies.userID) {
//     return <></>;
//   }

//   return (
//     <footer>
//       <div className="MerFooterItem" onClick={() => navigate("/")}>
//         <FaHome />
//         <p>Home</p>
//       </div>
//       <div className="MerFooterItem" onClick={() => navigate("/sell")}>
//         <FaCamera />
//         <p>Listing</p>
//       </div>
//       <div
//         className="MerFooterItem"
//         onClick={() => navigate(`/user/${cookies.userID}`)}
//       >
//         <FaUser />
//         <p>MyPage</p>
//       </div>
//     </footer>
//   );
// };
