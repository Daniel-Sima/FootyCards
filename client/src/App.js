import './App.css';
import React, { useCallback, useEffect, useState } from "react"
import { CssVarsProvider } from '@mui/joy/styles';
import theme from "./Theme"
import { Route, Routes } from "react-router-dom";
import Home from "./components/Home/Home"
import Login from './components/Login';
// import ErrorPage from './components/ErrorPage';
import SignUp from './components/SignUp';

function App() {
  const [jwtToken, setJwtToken] = useState("");
  const [tickInterval, setTickInterval] = useState();

  // Utilise avec le local storage si le cookies 
  // ne fonctionnent pas 


  const toggleRefresh = useCallback((status) => {
    if (status) {
      const requestOptions = {
        method: "GET",
        credentials: "include", // to have cookies
      }
      let i = setInterval(() => {
        fetch(`https://footycards-production-39e4.up.railway.app/refresh`, requestOptions)
          .then((response) => response.json())
          .then((data) => {
            if (data.access_token) {
              setJwtToken(data.access_token)
            }
          })
          .catch(error => {
            if (localStorage.getItem('jwtToken')) {
              setJwtToken(localStorage.getItem('jwtToken'));
            } else {
              setJwtToken("");
            }
          })
      }, 600_000) // 10 min 
      setTickInterval(i);
    } else {
      setTickInterval(null);
      clearInterval(tickInterval);
    }

  }, [tickInterval])

  useEffect(() => {
    if (jwtToken === "") {
      const requestOptions = {
        method: "GET",
        credentials: "include", // to have cookies
      }

      fetch(`https://footycards-production-39e4.up.railway.app/refresh`, requestOptions)
        .then((response) => response.json())
        .then((data) => {
          if (data.access_token) {
            setJwtToken(data.access_token)
            toggleRefresh(true);
          }
        })
        .catch(error => {
          if (localStorage.getItem('jwtToken')) {
            setJwtToken(localStorage.getItem('jwtToken'));
          } else {
            setJwtToken("");
          }
        })
    }
  }, [jwtToken, toggleRefresh])

  

  return (
    <CssVarsProvider theme={theme}>
      <Routes>
        <Route
          exact
          path="/*"
          element={jwtToken !== "" ?
            <Home setJwtToken={setJwtToken} toggleRefresh={toggleRefresh} sx={{ overflow: 'hidden' }} />
            :
            <Login setJwtToken={setJwtToken} toggleRefresh={toggleRefresh} />
          }
        >
        </Route>
        <Route
          exact
          path="/login"
          element={
            <Login setJwtToken={setJwtToken} toggleRefresh={toggleRefresh} />
          }
        >
        </Route>
        <Route
          path="/sign-up"
          element={<SignUp />}
        >
        </Route>
      </Routes>
    </CssVarsProvider>
  );
}

export default App;