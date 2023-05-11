import React, { useEffect, useState } from "react";
import PlayerCards from "../PlayerCards/PlayerCards";

const Team = (props) => {
  const [team, setTeam] = useState(null);
  const [club, setClub] = useState(null);

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
        let newTeam = []
        data.club.forEach(element => {
          if (element.inteam) newTeam.push(element)
        });
        setTeam(newTeam)
        setClub(data.club)
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

        fetch(`https://footycards-production.up.railway.app/userInfo`, requestOptions)
          .then((response) => response.json())
          .then((data) => {
            let newTeam = []
            data.club.forEach(element => {
              if (element.inteam) newTeam.push(element)
            });
            setTeam(newTeam)
            setClub(data.club)
          })
          .catch(err => {
            console.log(err);
          })
      })
  }, [props.user.team]);

  return (
    <>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        {team && club &&
          (club.length === 0 ?
            <h1 style={{ color: "white" }}>You don't have players yet, go to the Market and buy some!</h1>
            :
            (
              team.length !== 0 ?
                <PlayerCards players={team} user={props.user} setUser={props.setUser} fromClub={true} setTeam={setTeam} />
                :
                <h1 style={{ color: "white" }}>Go to the Club to define the players of your team! </h1>
            )
          )
        }
        {team && team.length < 4 &&
          <h1 style={{ color: "white" }}>
            You can add {4 - team.length} more players to maximize your score!
          </h1>
        }
        {team && team.length === 4 &&
          <h1 style={{ color: "white" }}>
            Your team is complete, you can wait to win score points or swap players if you want!
          </h1>
        }

        {!team &&
          (!club ?
            <h1 style={{ color: "white" }}>You don't have players yet, go to the Market and buy some!</h1>
            :
            <h1 style={{ color: "white" }}>Go to the Club to define the players of your team! </h1>
          )
        }
      </div>
    </>
  );
};

export default Team;
