import {FC, useEffect, useState} from "react";
import {Box, Card, CardContent, LinearProgress, Stack, Typography} from "@mui/material";
import CheckCircle from '@mui/icons-material/CheckCircle';
import axios from "axios";

interface Achievement {
    id: string;
    name: string;
    description: string;
    value: number;
    goal: number;
}

const App: FC = () => {
    const [achievements, setAchievements] = useState<Achievement[]>([]);

    const loadAchievements = async () => {
        try {
            const {data} = await axios.get(`http://localhost:8080/achievements`);
            setAchievements(data);
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        loadAchievements().catch(error => console.error(error));
    }, []);

    return (
        <Box display="flex" justifyContent="center" alignItems="center" flexWrap="wrap" gap={2} height="100vh"
             width="100vw" bgcolor="lightgrey">
            <Box position="absolute" top={0} display="flex" alignItems="center">
                <img src="/golang.png" alt="Logo" style={{height: 200}}/>
                <Typography variant="h4" component="div" gutterBottom color="text.primary">Your achievements in
                    GolangBank</Typography>
            </Box>
            {achievements.map(({id, name, description, value, goal}) => (
                <Card key={id} sx={{width: 300, height: 200}}>
                    <CardContent>
                        <Stack spacing={1}>
                            {value === goal ? (
                                <CheckCircle color="success" fontSize="large"/>
                            ) : (
                                <LinearProgress variant="determinate" value={(value / goal) * 100}/>
                            )}
                            <Typography variant="h6">{name}</Typography>
                            <Typography variant="body1">{description}</Typography>
                            <Typography variant="body2">Value: {value}</Typography>
                            <Typography variant="body2">Expected: {goal}</Typography>
                        </Stack>
                    </CardContent>
                </Card>
            ))}
        </Box>
    );
};

export default App;