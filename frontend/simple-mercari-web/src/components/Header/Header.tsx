
import "./Header.css";
import logo from '../assets/logo.png'
import * as React from 'react';
import { styled, alpha } from '@mui/material/styles'
import InputBase from '@mui/material/InputBase';
import SearchIcon from '@mui/icons-material/Search';
import {Toolbar} from "@mui/material";
import {useCookies} from "react-cookie";

export const Header: React.FC = () => {
    const [cookies] = useCookies(["userID", "token"]);

    const Search = styled('div')(({ theme }) => ({
        position: 'relative',
        borderRadius: theme.shape.borderRadius,
        backgroundColor: alpha(theme.palette.grey["300"], 0.15),
        border: '1px solid' + theme.palette.grey["300"],
        '&:hover': {
            border: '1px solid' + theme.palette.grey["400"],
        },
        marginRight: theme.spacing(2),
        marginLeft: 0,
        width: '100%',
        [theme.breakpoints.up('sm')]: {
            marginLeft: theme.spacing(2),
            width: 'auto',
        },
    }));

    const SearchIconWrapper = styled('div')(({ theme }) => ({
        padding: theme.spacing(0, 2),
        height: '100%',
        position: 'absolute',
        pointerEvents: 'none',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    }));

    const StyledInputBase = styled(InputBase)(({ theme }) => ({
        color: 'inherit',
        '& .MuiInputBase-input': {
            padding: theme.spacing(1, 1, 1, 0),
            // vertical padding + font size from searchIcon
            paddingLeft: `calc(1em + ${theme.spacing(4)})`,
            transition: theme.transitions.create('width'),
            width: '35ch',
            [theme.breakpoints.down('md')]: {
                width: '23ch',
            },
            [theme.breakpoints.down('sm')]: {
                width: '10ch',
                '&::placeholder': {
                    textOverflow: 'ellipsis !important',
                }
            },
        },
    }));

  return (
    <>
      <header>
          <Toolbar>
              <div>
                  <img className="logo" src={logo} alt="logo" />
              </div>
              {cookies.token &&
                  <Search>
                      <SearchIconWrapper>
                          <SearchIcon />
                      </SearchIconWrapper>
                      <StyledInputBase
                          placeholder="Shop one-of-a-kind finds"
                          inputProps={{ 'aria-label': 'search' }}
                      />
                  </Search>
              }
          </Toolbar>
      </header>
    </>
  );
}
