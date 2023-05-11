import React from "react";
import Typography from '@mui/joy/Typography';
import Sheet from '@mui/joy/Sheet';
import SportsSoccerIcon from '@mui/icons-material/SportsSoccer';
import Divider from '@mui/joy/Divider'
// import Avatar from '@mui/joy/Avatar'
import Box from '@mui/joy/Box';
import Button from '@mui/joy/Button';
// import IconButton from '@mui/joy/IconButton';
// import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import { Link, useNavigate } from "react-router-dom";



const Header = (props) => {

  const navigate = useNavigate();

  const logOut = () => {
    const requestOptions = {
      method: "GET",
      credentials: "include",
    }
    fetch(`/logout`, requestOptions)
    .catch(error => {
      console.log("error loggin out", error);
    })
    .finally(() => {
      props.toggleRefresh(false);
      props.setJwtToken("");
      localStorage.removeItem('jwtToken');
      localStorage.removeItem('email');
    })
    navigate("/login");
  }

  return (
    <Sheet
      variant="solid"
      color="primary"
      invertedColors
      sx={{
        width: "100vw",
        flexGrow: 1,
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        p: 2,
        minWidth: 'min-content'
      }}
    >
      <Box sx={{
        display: "flex",
        justifyContent: "left",
        alignItems: "center"
      }}>
        <SportsSoccerIcon sx={{ mt: 0.5, fontSize: 45, px: 1, color: "#FFFFFF" }}></SportsSoccerIcon>
        {/* <Link to="/" style={{ color: 'inherit', textDecoration: 'none' }}> */}
        <Typography level="h1" sx={{ fontSize: 30, pr: 2 }}> FootyCards </Typography>
        {/* </Link> */}
        <Divider orientation="vertical" color="primary" />
      </Box>

      <Box>
        <Link to="/club" style={{ color: 'inherit', textDecoration: 'none' }}>
          <Button color="primary" sx={{ mx: 5, p: 1.2, fontWeight: "bold", fontSize: 30, display: { xs: 'none', md: 'inline-flex' } }}>
            Club
          </Button>
        </Link>

        <Link to="/team" style={{ color: 'inherit', textDecoration: 'none' }}>
          <Button color="primary" sx={{ mx: 5, p: 1.2, fontWeight: "bold", fontSize: 30, display: { xs: 'none', md: 'inline-flex' } }}>
            Team
          </Button>
        </Link>

        <Link to="/market" style={{ color: 'inherit', textDecoration: 'none' }}>
          <Button color="primary" sx={{ mx: 5, p: 1.2, fontWeight: "bold", fontSize: 30, display: { xs: 'none', md: 'inline-flex' } }}>
            Market
          </Button>
        </Link>

        <Link to="/ranking" style={{ color: 'inherit', textDecoration: 'none' }}>
          <Button color="primary" sx={{ mx: 5, p: 1.2, fontWeight: "bold", fontSize: 30, display: { xs: 'none', md: 'inline-flex' } }}>
            Ranking
          </Button>
        </Link>
      </Box>

      {
        <Box>
          <Button variant="soft" sx={{ backgroundColor: '#b40d0d', color: 'white', '&:hover': { backgroundColor: '#b40d0d' } }}
            onClick={logOut}
          >
            Logout
          </Button>
        </Box>
      }

    </Sheet >

  );
}

export default Header;