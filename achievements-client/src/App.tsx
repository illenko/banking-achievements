import {FC, useEffect, useState} from "react";
import {Box, Card, CardContent, CircularProgress, Stack, Typography} from "@mui/material";
import CheckCircle from '@mui/icons-material/CheckCircle';
import axios from "axios";

interface Achievement {
    id: string;
    name: string;
    description: string;
    value: number;
    goal: number;
}

const fetchAchievements = async () => {
    try {
        const {data} = await axios.get(`http://localhost:8080/achievements`);
        return data;
    } catch (error) {
        console.error(error);
    }
};

const AchievementCard: FC<Achievement> = ({id, name, description, value, goal}) => (
    <Card key={id} sx={{width: 300, height: 220}}>
        <CardContent>
            <Stack spacing={1}>
                {value === goal ? (
                    <CheckCircle color="success" fontSize="large"/>
                ) : (
                    <CircularProgress variant="determinate" value={(value / goal) * 100}/>
                )}
                <Typography variant="h6">{name}</Typography>
                <Typography variant="body1">{description}</Typography>
                <Typography variant="body2"><b>Value:</b> {value}</Typography>
                <Typography variant="body2"><b>Goal:</b> {goal}</Typography>
            </Stack>
        </CardContent>
    </Card>
);

const App: FC = () => {
    const [achievements, setAchievements] = useState<Achievement[]>([]);

    useEffect(() => {
        fetchAchievements().then(setAchievements);
    }, []);

    return (
        <Box display="flex" flexDirection="column" height="100vh" width="100vw" bgcolor="lightgrey">
            <Box display="flex" alignItems="center" justifyContent="center" p={1} bgcolor="darkgrey">
                <img src="/golang.png" alt="Logo" style={{height: 50}}/>
                <Typography variant="h6" component="div" gutterBottom color="text.primary" ml={2}>
                    Your achievements in GolangBank
                </Typography>
            </Box>
            <Box display="flex" justifyContent="center" alignItems="center" flexWrap="wrap" gap={2} flexGrow={1}>
                {achievements.map(achievement => <AchievementCard {...achievement} />)}
            </Box>
        </Box>
    );
};

export default App;