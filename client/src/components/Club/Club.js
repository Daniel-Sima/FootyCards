import React, { useEffect, useState } from "react";
import PlayerCards from "../PlayerCards/PlayerCards";
import Select, { selectClasses } from '@mui/joy/Select';
import Option from '@mui/joy/Option';
import KeyboardArrowDown from '@mui/icons-material/KeyboardArrowDown';
import RadioButton from "./../Buttons/RadioButton";



const Club = (props) => {
    // const [club, setClub] = useState([]);
    const { user, setUser } = props;
    const [sortedDataMarket, setSortedDataMarket] = useState(null);
    const [availableClubs, setAvailableClubs] = useState(new Set());
    const [filteredDataMarket, setFilteredDataMarket] = useState([]);
    const [selectedValue, setSelectedValue] = useState(null); // pour enlever la valeur du placeHolder
    const fromClub = true;

    const handleEquipes = (club) => {
        setSelectedValue(null);
        if (club === "All") {
            setSortedDataMarket(user.club)
            setFilteredDataMarket(user.club)
        } else {
            const filteredData = user.club.filter(player => player.club.includes(club));
            setFilteredDataMarket(filteredData) // pour les equipes 
            setSortedDataMarket(filteredData) // sur les filtres
        }
    }

    const handleRatingGlobal = () => {
        const sorted = [...sortedDataMarket].sort((a, b) => b.ratingGlobal - a.ratingGlobal);
        setSortedDataMarket(sorted);
    }

    const handleRatingLast = () => {
        const sorted = [...sortedDataMarket].sort((a, b) => b.lastRating - a.lastRating);
        setSortedDataMarket(sorted);
    }

    const handlePrixCroissant = () => {
        const sorted = [...sortedDataMarket].sort((a, b) => a.price - b.price);
        setSortedDataMarket(sorted);
    }

    const handlePrixDecroissant = () => {
        const sorted = [...sortedDataMarket].sort((a, b) => b.price - a.price);
        setSortedDataMarket(sorted);
    }

    const handlePosition = (positionGet) => {
        const filteredData = filteredDataMarket.filter(player => player.position.includes(positionGet));
        setSortedDataMarket(filteredData)
    }

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`/userInfoCookies`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUser(data)
                // setClub(data.club)
                data.club && setSortedDataMarket(data.club)
                data.club && setFilteredDataMarket(data.club)
                data.club && setAvailableClubs((new Set(data.club.map((item) => item.club))).add("All"));
            })
            .catch(err => {
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

                fetch(`/userInfo`, requestOptions)
                    .then((response) => response.json())
                    .then((data) => {
                        setUser(data)
                        setSortedDataMarket(data.club)
                        setFilteredDataMarket(data.club)
                        setAvailableClubs((new Set(data.club.map((item) => item.club))).add("All"));
                    })
                    .catch(err => {
                        console.log(err);
                    })
            })
    }, [setUser]); // setUser mby

    return (
        <>
            {((user.club && user.club.length === 0) || (!user.club)) &&
                <div
                    style={{
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                        alignItems: "center",
                        height: "100vh",
                    }}
                >

                    <h1 style={{ color: "white" }}>You don't have players yet, go to the Market and buy some!</h1>
                </div>
            }

            <Select
                style={{ fontSize: '16px', position: 'relative', top: '0px', left: '750px', fontFamily: 'Impact, sans-serif', fontWeight: 'bold' }}
                color="primary"
                variant="solid"
                size="sm"
                placeholder="Sorted by ..."
                indicator={<KeyboardArrowDown />}
                value={selectedValue}
                sx={{
                    width: 215,
                    [`& .${selectClasses.indicator}`]: {
                        transition: '0.2s',
                        [`&.${selectClasses.expanded}`]: {
                            transform: 'rotate(-180deg)',
                        },
                    },
                }}
            >
                <Option onClick={handlePrixCroissant} value="ascending price">Ascending price</Option>
                <Option onClick={handlePrixDecroissant} value="decreasing price">Decreasing price</Option>
                <Option onClick={handleRatingGlobal} value="rating global">Rating global</Option>
                <Option onClick={handleRatingLast} value="last global">Last global</Option>
                <Option onClick={() => { handlePosition("Attacker") }} value="attacker">Attacker</Option>
                <Option onClick={() => { handlePosition("Midfielder") }} value="midfielders">Midfielders</Option>
                <Option onClick={() => { handlePosition("Defender") }} value="defenders">Defenders</Option>
                <Option onClick={() => { handlePosition("Goalkeeper") }} value="goalkeeper">Goalkeepers</Option>
            </Select>

            <RadioButton availableClubs={availableClubs} handleEquipes={handleEquipes}></RadioButton>

            {sortedDataMarket && <PlayerCards user={props.user} setUser={props.setUser} sortedMarket={sortedDataMarket} setMarket={setSortedDataMarket} fromClub={fromClub} />}

            <hr style={{ height: '100px' }} />
        </>
    )
}

export default Club;