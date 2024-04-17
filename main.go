package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Name   string
	Gender string
	Age    int
}

type Vehicle struct {
	Owner       string
	Model       string
	NumberPlate string
}

type Ride struct {
	ID             int
	Driver         string
	Origin         string
	Destination    string
	AvailableSeats int
	Vehicle        Vehicle
	Active         bool
}

var nameUserMap = make(map[string]User)
var ownerVehiclesMap = make(map[string][]Vehicle)
var rides = make(map[int]Ride, 0)
var ownerRideStats = make(map[string]map[string]int)

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func addUser(userDetail string) {
	userDetails := strings.Split(userDetail, ", ")
	if len(userDetails) != 3 {
		fmt.Println("Invalid user details")
		return
	}
	user := User{Name: userDetails[0], Gender: userDetails[1], Age: toInt(userDetails[2])}
	nameUserMap[user.Name] = user
	ownerRideStats[user.Name] = map[string]int{"offered": 0, "taken": 0}
	fmt.Println("User added:", user)
}

func addVehicle(vehicleDetail string) {
	vehicleDetails := strings.Split(vehicleDetail, ", ")
	if len(vehicleDetails) != 3 {
		fmt.Println("Invalid vehicle details")
		return
	}
	vehicleData := Vehicle{Owner: vehicleDetails[0], Model: vehicleDetails[1], NumberPlate: vehicleDetails[2]}
	ownerVehiclesMap[vehicleData.Owner] = append(ownerVehiclesMap[vehicleData.Owner], vehicleData)
	fmt.Println("Vehicle added for", vehicleData.Owner, ":", vehicleData)
}

func offerRide(rideDetail string) {
	rideDetails := strings.Split(rideDetail, ", ")
	if len(rideDetails) != 6 {
		fmt.Println("Invalid vehicle details")
		return
	}
	driver := rideDetails[0]
	origin := strings.Split(rideDetails[1], "=")[1]
	availableSeats := toInt(strings.Split(rideDetails[2], "=")[1])
	vehicleModel := strings.Split(rideDetails[3], "=")[1]
	numberPlate := rideDetails[4]
	destination := strings.Split(rideDetails[5], "=")[1]

	vehicleData := findVehicle(driver, vehicleModel, numberPlate)
	if vehicleData == nil {
		fmt.Println("Vehicle not found")
		return
	}

	for _, r := range rides {
		if r.Vehicle == *vehicleData && r.Active {
			fmt.Println("Ride already active for this vehicle")
			return
		}
	}

	ride := Ride{
		ID:             len(rides),
		Driver:         driver,
		Origin:         origin,
		Destination:    destination,
		AvailableSeats: availableSeats,
		Vehicle:        *vehicleData,
		Active:         true,
	}
	rides[len(rides)] = ride
	ownerRideStats[driver]["offered"]++
	fmt.Println("Ride offered:: and ride details:: ", ride)
}

func findVehicle(owner, vehicleModel, numberPlate string) *Vehicle {
	for _, v := range ownerVehiclesMap[owner] {
		if vehicleModel == v.Model && numberPlate == v.NumberPlate {
			return &v
		}
	}
	return nil
}

func selectRide(user, source, destination string, seats int, selectionStrategy string) (int, error) {
	filteredRides := make([]Ride, 0)
	for _, ride := range rides {
		if ride.Origin == source && ride.Destination == destination && ride.AvailableSeats >= seats && ride.Active {
			filteredRides = append(filteredRides, ride)
		}
	}

	if len(filteredRides) == 0 {
		return -1, errors.New("no rides found")
	}

	selectedRide := filteredRides[0]
	if selectionStrategy == "Most Vacant" {
		maxSeats := 0
		for _, r := range filteredRides {
			if r.AvailableSeats > maxSeats {
				selectedRide = r
				maxSeats = r.AvailableSeats
			}
		}
	} else if strings.HasPrefix(selectionStrategy, "Preferred Vehicle=") {
		preferredVehicle := strings.Split(selectionStrategy, "=")[1]
		found := false
		for _, r := range filteredRides {
			if r.Vehicle.Model == preferredVehicle {
				selectedRide = r
				found = true
				break
			}
		}
		if !found {
			return -1, errors.New("no rides found with the preferred vehicle")
		}
	}

	if ride, ok := rides[selectedRide.ID]; ok {
		ride.AvailableSeats -= seats
		rides[ride.ID] = ride
	}

	if _, ok := ownerRideStats[user]; !ok {
		ownerRideStats[user]["taken"]++
	} else {
		ownerRideStats[user] = map[string]int{"offered": 0, "taken": 0}
		ownerRideStats[user]["taken"]++
	}
	fmt.Println("Ride selected:: ", selectedRide)
	return selectedRide.ID, nil
}

func endRide(rideID int) {
	if rideID < 0 || rideID >= len(rides) {
		fmt.Println("Invalid ride ID")
		return
	}
	ride := rides[rideID]
	ride.Active = false
	rides[ride.ID] = ride
	fmt.Println("Ride ended:", rides[rideID])
}

func printExistingRides() {
	fmt.Println("=========EXISTING RIDES=============")
	for _, ride := range rides {
		fmt.Printf("ID: %d -- Driver: %s -- Origin: %s -- Destination: %s -- AvailableSeats: %d -- Vehicle: %s -- Active:%v\n", ride.ID, ride.Driver, ride.Origin, ride.Destination, ride.AvailableSeats, ride.Vehicle, ride.Active)
	}
}

func printRideStats() {
	fmt.Printf("==============RIDE STATS================\n")
	for user, stats := range ownerRideStats {
		fmt.Printf("%s: %d Taken, %d Offered\n", user, stats["taken"], stats["offered"])
	}
	fmt.Printf("========================================\n")
}

func toInt(input string) int {
	var value int
	fmt.Sscanf(input, "%d", &value)
	return value
}

func main() {
	fmt.Println("POPULATING DUMMY DATA FOR USER, VEHICLE and OFFER")
	addDummyData()
	for {
		fmt.Printf("==============MENU================\n")
		fmt.Println("1. Add New User Details")
		fmt.Println("2. Add New Vehicle Details")
		fmt.Println("3. Offer Ride")
		fmt.Println("4. Select Ride")
		fmt.Println("5. End Ride")
		fmt.Println("6. Print Ride Stats")
		fmt.Println("7. Print Existing Rides")
		fmt.Println("8. Quit")
		fmt.Printf("==================================\n\n")
		choice := readInput("Enter your choice: ")

		switch choice {
		case "1":
			userDetail := readInput("Enter user details in Format :: [Name, Gender, Age] : ")
			addUser(userDetail)
		case "2":
			vehicleDetail := readInput("Enter vehicle details in Format :: [Owner, Model, License Plate] : ")
			addVehicle(vehicleDetail)
		case "3":
			rideDetail := readInput("Enter ride details ([Driver_Name], Origin=..., Available Seats=..., Vehicle=..., [NumberPlate] Destination=...): ")
			offerRide(rideDetail)
		case "4":
			user := readInput("Enter User:")
			source := readInput("Enter source: ")
			destination := readInput("Enter destination: ")
			seats := readInput("Enter number of seats: ")
			strategy := readInput("Enter selection strategy (Most Vacant/Preferred Vehicle=...): ")
			selectedRideID, _ := selectRide(user, source, destination, toInt(seats), strategy)
			fmt.Printf("Selected Ride id is:: %d\n", selectedRideID)
		case "5":
			rideID := toInt(readInput("Enter ride ID to end: "))
			endRide(rideID)
		case "6":
			printRideStats()
		case "7":
			printExistingRides()
		case "8":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func addDummyData() {
	addUser("Rohan, M, 36")
	addVehicle("Rohan, Swift, KA-01-12345")
	addUser("Shashank, M, 29")
	addVehicle("Shashank, Baleno, TS-05-62395")
	addUser("Nandini, F, 29")
	addUser("Shipra, F, 27")
	addVehicle("Shipra, Polo, KA-05-41491")
	addVehicle("Shipra, Activa, KA-12-12332")
	addUser("Gaurav, M, 29")
	addUser("Rahul, M, 35")
	addVehicle("Rahul, XUV, KA-05-1234")
	addVehicle("Rohan, Polo, KA-01-44252")

	offerRide("Rohan, Origin=Hyderabad, Available Seats=1, Vehicle=Swift, KA-01-12345, Destination=Bangalore")
	offerRide("Shipra, Origin=Bangalore, Available Seats=1, Vehicle=Activa, KA-12-12332, Destination=Mysore")
	offerRide("Shipra, Origin=Bangalore, Available Seats=2, Vehicle=Polo, KA-05-41491, Destination=Mysore")
	offerRide("Shashank, Origin=Hyderabad, Available Seats=2, Vehicle=Baleno, TS-05-62395, Destination=Bangalore")
	offerRide("Rahul, Origin=Pune, Available Seats=5, Vehicle=XUV, KA-05-1234, Destination=Bangalore")
	offerRide("Rohan, Origin=Mumbai, Available Seats=1, Vehicle=Swift, KA-01-12345, Destination=Delhi")
	offerRide("Rohan, Origin=Mumbai, Available Seats=1, Vehicle=Polo, KA-01-44252, Destination=Pune")

	selectedRideID, _ := selectRide("Nandini", "Bangalore", "Mysore", 1, "Most Vacant")
	fmt.Printf("Selected Ride id is:: %d\n", selectedRideID)
	selectedRideID, _ = selectRide("Gaurav", "Bangalore", "Mysore", 1, "Preferred Vehicle=Activa")
	fmt.Printf("Selected Ride id is:: %d\n", selectedRideID)
	selectedRideID, _ = selectRide("Shashank", "Mumbai", "Mysore", 1, "Most Vacant")
	if selectedRideID == -1 {
		multiRides := findPossibleRides("Mumbai", "Mysore", 1)
		fmt.Println("POSSIBLE RIDES")
		if len(multiRides) > 0 {
			for _, ride := range multiRides[0] {
				fmt.Println(ride)
			}
		}
	}
	fmt.Printf("Selected Ride id is:: %d\n", selectedRideID)
	selectedRideID, _ = selectRide("Rohan", "Hyderabad", "Bangalore", 1, "Preferred Vehicle=Baleno")
	fmt.Printf("Selected Ride id is:: %d\n", selectedRideID)
	selectedRideID, _ = selectRide("Shashank", "Hyderabad", "Bangalore", 1, "Preferred Vehicle=Polo")
	fmt.Printf("Selected Ride id is:: %d\n", selectedRideID)

	printRideStats()
	printExistingRides()
}

func findPossibleRides(source, destination string, seats int) [][]Ride {
	var possibleRoutes [][]Ride

	// Recursive function to perform depth-first search
	var dfs func(string, []Ride)
	dfs = func(currentLocation string, path []Ride) {
		for _, ride := range rides {
			if ride.Origin == currentLocation && ride.Active && seats <= ride.AvailableSeats {
				// Check if the ride reaches the destination
				if ride.Destination == destination && ride.Active && seats <= ride.AvailableSeats {
					// If the ride reaches the destination, add it to the current path
					newPath := make([]Ride, len(path))
					copy(newPath, path)
					newPath = append(newPath, ride)
					possibleRoutes = append(possibleRoutes, newPath)
				} else {
					// Continue exploring from the destination of the current ride
					newPath := make([]Ride, len(path))
					copy(newPath, path)
					newPath = append(newPath, ride)
					dfs(ride.Destination, newPath)
				}
			}
		}
	}

	// Start DFS from the source location
	dfs(source, []Ride{})

	return possibleRoutes
}
