import React from "react";
import PlayerCard from "../PlayerCard/PlayerCard";

const PlayerCards = (props) => {
  return (
    <>
    {
    props.players ? 
      <div style={{ display: "flex", flexWrap: "wrap" }}>
      {props.players &&
        props.players.map((item, index) => (
          <PlayerCard
            key={item._id}
            data={item}
            team={props.players}
            user={props.user}
            setUser={props.setUser}
            fromClub={props.fromClub}
            setTeam={props.setTeam}
          />
        ))}
    </div>
    :
      <div style={{ display: "flex", flexWrap: "wrap" }}>
        {props.sortedMarket &&
          props.sortedMarket.map((item, index) => (
            <PlayerCard
              key={item._id}
              data={item}
              user={props.user}
              setUser={props.setUser}
              fromClub={props.fromClub}
              sortedMarket = {props.sortedMarket}
              setMarket={props.setMarket}
            />
          ))}
      </div>
    }

    </>
  );
};

export default PlayerCards;
