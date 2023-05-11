import React, { useEffect, useState } from "react"
import { Routes, Route } from 'react-router-dom'
// import Header from "../Header/Header"
import Club from "./../Club/Club"
// import { Container } from "@mui/system";
// import { CssBaseline } from "@mui/material";
import Box from '@mui/joy/Box';
// import StadiumImage from "./../../assets/stadiumBG.jpg"
// import Mbappe from "./../../assets/mbappe.jpg"
// import { useState } from "react";
import Header from "./../Header/Header"
import Team from "./../Team/Team"
import Market from "./../Market"
import Ranking from "./../Ranking/Ranking"


const Home = (props) => {

    const [user, setUser] = useState([]);

    const handleConnection = () => {
        let payload = {
            email: localStorage.getItem('email'),
        }

        const requestOptions = {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload),
        }

        fetch(`https://footycards-production.up.railway.app/userInfo`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUser(data)
            })
            .catch(err => {
                console.log(err);
            })
    }

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");


        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`https://footycards-production.up.railway.app/userInfoCookies`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUser(data)
            })
            .catch(err => {
                handleConnection();
            })

    }, []); // setUser mby

    return (
        <div>
            <Header setJwtToken={props.setJwtToken} toggleRefresh={props.toggleRefresh} />

            <Box
                sx={{
                    backgroundImage: `linear-gradient(0deg, #334ec1 0%, #000a19 60%)`,
                    backgroundPosition: 'center',
                    height: "90dvh",
                    overflow: "auto"
                }}
            >
                <Routes>
                    <Route
                        exact
                        path="/club"
                        element={<Club user={user} setUser={setUser} />}
                    >
                    </Route>
                    <Route
                        path="/team"
                        element={<Team user={user} setUser={setUser} />}
                    >
                    </Route>
                    <Route
                        path="/market"
                        element={<Market user={user} setUser={setUser} />}
                    >
                    </Route>
                    <Route
                        path="/ranking"
                        element={<Ranking user={user} />}
                    >
                    </Route>
                </Routes>
            </Box>
        </div>
    );

}

export default Home;