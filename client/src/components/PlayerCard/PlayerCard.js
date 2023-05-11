import React, { useState } from "react";
import Card from "@mui/joy/Card";
import CardOverflow from "@mui/joy/CardOverflow";
import AspectRatio from "@mui/joy/AspectRatio";
import Typography from "@mui/joy/Typography";
import Button from "@mui/joy/Button";

const PlayerCard = (props) => {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const addPlayerTeam = (player, userID) => {
    const requestOptions = {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        userID: userID,
        player: player,
      }),
    };
    fetch(`team`, requestOptions)
      .then(async (response) => { // pq async ?
        const status = await response.status;
        if (status === 200) {
          const newPlay = {
            ...player,
            inteam: true,
          };
          const newMarket = [...props.sortedMarket]
          let index = newMarket.indexOf(player)
          if (index !== -1) {
            newMarket[index] = newPlay
          }
          props.setMarket(newMarket)
        } else if (status === 400) {
          alert("You already have max player in your team (4)");
        }
      })
      .then((data) => data)
      .catch((err) => {
        console.log(err);
      });
  };

  const deletePlayerTeam = (player, userID) => {
    const requestOptions = {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        userID: userID,
        player: player,
      }),
    };
    fetch(`team`, requestOptions)
      .then(async (response) => {
        const status = await response.status;
        console.log(status)
        if (status === 200) {
          let newTeam = [...props.team]
          let indexTeam = newTeam.indexOf(player)
          if (indexTeam > -1) {
            newTeam.splice(indexTeam, 1);
          }
          props.setTeam(newTeam)
        }
      })
      .then((data) => data)
      .catch((err) => {
        console.log(err);
      });
  };

  const deletePlayerTeamFromClub = (player, userID) => {
    const requestOptions = {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        userID: userID,
        player: player,
      }),
    };
    fetch(`team`, requestOptions)
      .then(async (response) => {
        const status = await response.status;
        // console.log(status)
        if (status === 200) {
          const newPlay = {
            ...player,
            inteam: false,
          };
          const newMarket = [...props.sortedMarket]
          let index = newMarket.indexOf(player)
          if (index !== -1) {
            newMarket[index] = newPlay
          }
          props.setMarket(newMarket)
        }
      })
      .then((data) => data)
      .catch((err) => {
        console.log(err);
      });
  };

  // Pour la MAJ des coins restants
  const getUserInfos = () => {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const requestOptions = {
      method: "GET",
      headers: headers,
    };

    fetch(`/userInfoCookies`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        props.setUser(data);
      })
      .catch((err) => {
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
            props.setUser(data)
          })
          .catch(err => {
            console.log(err);
          })
      });
  };

  const getMarket = () => {
    console.log("Get Market after buying exemplar")
    const headers = new Headers();
    headers.append("Content-Type", "application/json");


    const requestOptions = {
      method: "GET",
      headers: headers,
    }

    fetch(`/allPlayers`, requestOptions)
      .then((response) => response.json())
      .then((dataMarket) => {
        props.setSortedDataMarket(dataMarket);
      })
      .catch(err => {
        console.log(err);
      })
  }

  const handleBuy = () => {
    const confirmed = window.confirm("Are you sure you want this player?");
    if (confirmed) {
      setIsSubmitting(true);

      const requestOptions = {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          user: props.user,
          player: props.data,
        }),
      };

      fetch(`/buyingPlayer`, requestOptions)
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            alert("[FAILED] " + data.message);
          } else {
            getMarket();
            getUserInfos();
            alert("[SUCCESS] You have well buyed this player for " + props.data.price + "!");
            props.data.exemplars = props.data.exemplars - 1
          }
        })
        .catch((err) => {
          alert(err);
        })
        .finally(() => {
          setIsSubmitting(false);
        });
    }
  };

  const handleSell = () => {
    const confirmed = window.confirm(
      "Are you sure you want to sell this player?"
    );
    if (confirmed) {
      setIsSubmitting(true);
      const requestOptions = {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
          user: props.user,
          player: props.data
        }),
      }

      fetch(`/sellingPlayer`, requestOptions)
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            alert("[FAILED] " + data.message);
          } else {
            const newMarket = [...props.sortedMarket].filter(player => player._id !== props.data._id);
            props.setMarket(newMarket);
            getUserInfos()
            alert("[SUCCESS] You sold this player well for " + props.data.price + "!")
          }
        })
        .catch(err => {
          alert(err);
        })
        .finally(() => {
          setIsSubmitting(false);
        })
    }
  };

  return (
    <>
      <Card color="primary" sx={{ width: 300, height: "auto", margin: "50px" }}>
        <CardOverflow sx={{ padding: 4, paddingLeft: 5, paddingRight: 5 }}>
          <AspectRatio ratio="1" sx={{ borderRadius: "10px" }}>
            <img src={props.data.image} alt="playerImg" />
          </AspectRatio>
        </CardOverflow>
        <Typography level="h5" sx={{ mt: -3, fontWeight: "bold" }}>
          {props.data.name}
        </Typography>
        <Typography level="body1" sx={{ mt: 0.5, mb: 2, fontWeight: "bold" }}>
          {props.data.club} | {props.data.position}
        </Typography>
        <Typography level="body1" sx={{ mt: -2, fontWeight: "bold" }}>
          Rating global: {props.data.ratingGlobal}
        </Typography>
        <Typography
          level="body1"
          sx={{ mt: 0, fontWeight: "bold", color: "#660000" }}
        >
          Last rating: {props.data.lastRating}
        </Typography>
        {props.fromMarket && (
          <>
            <Typography
              level="body1"
              sx={{
                ml: 19,
                mt: -0.25,
                mb: 2,
                fontWeight: "bold",
                border: "3px solid black",
              }}
            >
              {props.data.exemplars} exemplars
            </Typography>
            <Button
              disabled={isSubmitting}
              sx={{ fontSize: "md", mt: 0 }}
              size="md"
              variant="solid"
              color="success"
              onClick={handleBuy}
            >
              {isSubmitting ? "Buying.. " : "Buy for " + props.data.price}
            </Button>
          </>
        )}

        {props.fromClub && (!props.data.inteam) &&
          <Button
            onClick={() => {
              addPlayerTeam(props.data, props.user._id);
            }}
          >
            Add player to team
          </Button>
        }

        {props.fromClub && props.data.inteam && !props.sortedMarket &&
          <Button
            sx={{
              backgroundColor: "#b40d0d",
              color: "white",
              "&:hover": { backgroundColor: "#b40d0d" },
            }}
            onClick={() => {
              deletePlayerTeam(props.data, props.user._id);
            }}
          >
            Remove player from team
          </Button >
        }

        {props.fromClub && props.data.inteam && props.sortedMarket &&
          <Button
            sx={{
              backgroundColor: "#b40d0d",
              color: "white",
              "&:hover": { backgroundColor: "#b40d0d" },
            }}
            onClick={() => {
              deletePlayerTeamFromClub(props.data, props.user._id);
            }}
          >
            Remove player from team
          </Button >
        }

        {props.fromClub && !props.data.inteam &&
          <Button
            disabled={isSubmitting}
            sx={{ fontSize: "md", mt: 1 }}
            size="md"
            variant="solid"
            color="danger"
            onClick={handleSell}
          >
            {isSubmitting ? "Selling.. " : "Sell for " + props.data.price}
          </Button>
        }
      </Card >
    </>
  );
};

export default PlayerCard;
