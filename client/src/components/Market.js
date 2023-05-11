import React, { useEffect, useState } from "react";

import PlayerCard from "./PlayerCard/PlayerCard"
import Chip from '@mui/joy/Chip';
import SavingsIcon from '@mui/icons-material/Savings';
import Select, { selectClasses } from '@mui/joy/Select';
import Option from '@mui/joy/Option';
import KeyboardArrowDown from '@mui/icons-material/KeyboardArrowDown';
import RadioButton from "./Buttons/RadioButton";

const Market = (props) => {
    const [dataMarket, setDataMarket] = useState([]); // All players
    const [sortedDataMarket, setSortedDataMarket] = useState(dataMarket);
    const [filteredDataMarket, setFilteredDataMarket] = useState(dataMarket);
    const [availableClubs, setAvailableClubs] = useState(new Set());
    const [selectedValue, setSelectedValue] = useState(null); // pour enlever la valeur du placeHolder
    const fromMarket = true;

    const handleEquipes = (club) => {
        setSelectedValue(null);
        if (club === "All") {
            setSortedDataMarket(dataMarket)
            setFilteredDataMarket(dataMarket) // pour les equipes 
        } else {
            const filteredData = dataMarket.filter(player => player.club.includes(club));
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

    const handleAscendingExemplars = () => {
        const sorted = [...sortedDataMarket].sort((a, b) => a.exemplars - b.exemplars);
        setSortedDataMarket(sorted);
    }

    const handleDecreasingExemplars = () => {
        const sorted = [...sortedDataMarket].sort((a, b) => b.exemplars - a.exemplars);
        setSortedDataMarket(sorted);
    }

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");


        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`https://footycards-production.up.railway.app/allPlayers`, requestOptions)
            .then((response) => response.json())
            .then((dataMarket) => {
                setDataMarket(dataMarket)
                setSortedDataMarket(dataMarket);
                setFilteredDataMarket(dataMarket);
                setAvailableClubs((new Set(dataMarket.map((item) => item.club))).add("All"));
            })
            .catch(err => {
                console.log(err);
            })
    }, []);

    return (
        <>
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
                <Option onClick={handleAscendingExemplars}  value="ascending exemplars">Ascending exemplars</Option>
                <Option onClick={handleDecreasingExemplars} value="decreasing exemplars">Decreasing exemplars</Option>
                <Option onClick={() => {handlePosition("Attacker")}}  value="attacker">Attacker</Option>
                <Option onClick={() => {handlePosition("Midfielder")}} value="midfielders">Midfielders</Option>
                <Option onClick={() => {handlePosition("Defender")}} value="defenders">Defenders</Option>
                <Option onClick={() => {handlePosition("Goalkeeper")}}value="goalkeeper">Goalkeepers</Option>
            </Select>

            <RadioButton availableClubs={availableClubs} handleEquipes={handleEquipes}></RadioButton>

            <Chip style={{ fontSize: '24px', position: 'fixed', top: '75px', right: '10px', fontFamily: 'Impact, sans-serif', fontWeight: 'bold' }} variant="solid" endDecorator={<SavingsIcon style={{ fontSize: '1.2em' }} />} >
                You have {props.user && props.user.coins} coins</Chip>

            <div style={{ display: "flex", flexWrap: "wrap" }}>
                {sortedDataMarket && sortedDataMarket.map((item, index) => (
                    <PlayerCard key={item._id} data={item} user={props.user} setUser={props.setUser} fromMarket={fromMarket} setSortedDataMarket={setSortedDataMarket} />
                ))}
            </div>
            <hr style={{height: '100px'}}/>
        </>
    )
}

export default Market;