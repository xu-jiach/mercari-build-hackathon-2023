
import "./Header.css";
import logo from '../assets/logo.png'
import * as React from 'react';
import { styled, alpha } from '@mui/material/styles'
import InputBase from '@mui/material/InputBase';
import SearchIcon from '@mui/icons-material/Search';
import {Toolbar} from "@mui/material";
import {useCookies} from "react-cookie";
import {useNavigate} from "react-router-dom";
import Avatar from '@mui/material/Avatar';
import Stack from '@mui/material/Stack';

// auto generates icons based on user name


function stringAvatar(character: string) {
  const name = character.length > 1 ? character : `${character}${character}`;
  const colors = ['#06CBFF', '#FF5757'];
  const randomIndex = Math.floor(Math.random() * colors.length);
  const color = colors[randomIndex];

  const avatarSx = {
    bgcolor: color,
    cursor: 'pointer', // Add the cursor property here
  };

  return {
    sx: avatarSx,
    children: name,
  };
}




export const Header: React.FC = () => {
    const [cookies] = useCookies(["userID", "token"])
    const navigate = useNavigate();

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

    const handleSubmit = (
        e: React.MouseEvent<HTMLButtonElement> | React.KeyboardEvent<HTMLInputElement>, search: string
    ) => {
        e.preventDefault();
        navigate(`/search-advanced?keyword=${search}`);
        window.location.reload();
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.nativeEvent.isComposing || e.key !== 'Enter') return
        const search = e.currentTarget.value;
        handleSubmit(e, search)
    }

  return (
    <>
      <header>
          <Toolbar>
              <div>
                  <img className="logo" src={logo} alt="logo" onClick={() => navigate('/')} />
              </div>
              {cookies.token &&
                <>
                  <Search>
                    <SearchIconWrapper>
                        <SearchIcon />
                    </SearchIconWrapper>
                    <StyledInputBase
                        placeholder="Shop one-of-a-kind finds"
                        inputProps={{ 'aria-label': 'search' }}
                        onKeyDown={handleKeyDown}
                        // TODO: Display suggestions when user types
                        // onChange={onChangeSuggestion}
                    />
                  </Search>
                </>
              }
          </Toolbar>
          {cookies.token &&
            <Stack sx={{ pr: 5, mt: 1.5 }}>
              <Avatar {...stringAvatar(cookies.userID)}
                onClick={() => navigate(`/user/${cookies.userID}`)}
              />
            </Stack>
          }
      </header>
    </>
  );
}
