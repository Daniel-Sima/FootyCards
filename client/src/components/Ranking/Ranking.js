import React, { useEffect, useState } from "react";
import Table from '@mui/joy/Table'

const Ranking = (props) => {

    const [ranking, setRanking] = useState([])

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");


        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`https://footycards-production-39e4.up.railway.app/ranking`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setRanking(data)
            })
            .catch(err => {
                console.log(err);
            })

    }, []);

    return (
        <>
            <Table
                borderAxis="xBetween"
                color="neutral"
                size="md"
                stickyHeader
                variant="plain"
                sx={{ padding: "25px" }}
            >
                <thead>
                    <tr>
                        <th style={{ width: '30%' }}>Rank</th>
                        <th>Score</th>
                        <th>Pseudo</th>
                        <th>Club name</th>
                    </tr>
                </thead>
                <tbody>
                    {ranking ? ranking.map((item, index) => (
                        <tr>
                            <td style={{ color: props.user.email === item.email ? "red" : "white" }}>{index+1} {/* {item.classement}*/ }</td>   {/* car deja tries */}
                            <td style={{ color: props.user.email === item.email ? "red" : "white" }}>{item.score}</td> 
                            <td style={{ color: props.user.email === item.email ? "red" : "white" }}>{item.pseudo}</td>
                            <td style={{ color: props.user.email === item.email ? "red" : "white" }}>{item.clubName}</td>
                        </tr>)) : <></>}
                </tbody>
            </Table>
            {/* <div style={{display: "flex", flexWrap:"wrap"}}> */}

            {/* </div> */}
        </>
        // <h1 style={{ position: 'absolute', top: '100px', left: '0', color: 'white' }}>Ranking</h1>
    )
}

export default Ranking;