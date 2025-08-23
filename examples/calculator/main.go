package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/orientallines/gossie/pkg/gossie"
)

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Advanced Calculator",
		AppDescription: "A powerful calculator with multiple operations and features",
		AppVersion:     "2.0.0",
		AppAuthor:      "Gossie Calculator Team",
	})

	// Basic arithmetic operations
	app.Command("add", func(cmd *gossie.Command) {
		cmd.Description("Add two or more numbers")
		cmd.Arg("numbers").Multiple().Description("Numbers to add")
		cmd.Flag("precision", "Decimal precision").Short('p')
		cmd.Action(func(c *gossie.Context) error {
			args := c.Args()
			if len(args) < 2 {
				return fmt.Errorf("at least 2 numbers required")
			}

			var sum float64
			for _, arg := range args {
				num, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return fmt.Errorf("invalid number: %s", arg)
				}
				sum += num
			}

			precision := 2
			if precStr := c.GetFlag("precision"); precStr != "" {
				if p, err := strconv.Atoi(precStr); err == nil && p >= 0 && p <= 10 {
					precision = p
				}
			}

			fmt.Printf("Sum: %.*f\n", precision, sum)
			return nil
		})
	})

	// Subtraction with multiple numbers
	app.Command("subtract", func(cmd *gossie.Command) {
		cmd.Description("Subtract numbers from the first number")
		cmd.Arg("start").Required().Description("Starting number")
		cmd.Arg("subtractors").Multiple().Description("Numbers to subtract")
		cmd.Action(func(c *gossie.Context) error {
			args := c.Args()
			if len(args) < 2 {
				return fmt.Errorf("at least 2 numbers required")
			}

			start, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid starting number: %s", args[0])
			}

			result := start
			for i := 1; i < len(args); i++ {
				num, err := strconv.ParseFloat(args[i], 64)
				if err != nil {
					return fmt.Errorf("invalid number: %s", args[i])
				}
				result -= num
			}

			fmt.Printf("%.2f - %s = %.2f\n", start, args[1:], result)
			return nil
		})
	})

	// Multiplication
	app.Command("multiply", func(cmd *gossie.Command) {
		cmd.Description("Multiply numbers together")
		cmd.Arg("factors").Multiple().Description("Numbers to multiply")
		cmd.Flag("show-steps", "Show multiplication steps")
		cmd.Action(func(c *gossie.Context) error {
			args := c.Args()
			if len(args) < 2 {
				return fmt.Errorf("at least 2 numbers required")
			}

			result := 1.0
			var numbers []float64

			for _, arg := range args {
				num, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return fmt.Errorf("invalid number: %s", arg)
				}
				numbers = append(numbers, num)
				result *= num
			}

			if c.HasFlag("show-steps") {
				fmt.Print(numbers[0])
				for i := 1; i < len(numbers); i++ {
					fmt.Printf(" × %.2f", numbers[i])
				}
				fmt.Printf(" = %.2f\n", result)
			} else {
				fmt.Printf("Product: %.2f\n", result)
			}

			return nil
		})
	})

	// Division with remainder
	app.Command("divide", func(cmd *gossie.Command) {
		cmd.Description("Divide two numbers")
		cmd.Arg("dividend").Required().Description("Number to divide")
		cmd.Arg("divisor").Required().Description("Number to divide by")
		cmd.Flag("remainder", "Show remainder")
		cmd.Action(func(c *gossie.Context) error {
			dividend, err := strconv.ParseFloat(c.GetArg("dividend"), 64)
			if err != nil {
				return fmt.Errorf("invalid dividend: %s", c.GetArg("dividend"))
			}

			divisor, err := strconv.ParseFloat(c.GetArg("divisor"), 64)
			if err != nil {
				return fmt.Errorf("invalid divisor: %s", c.GetArg("divisor"))
			}

			if divisor == 0 {
				return fmt.Errorf("division by zero")
			}

			quotient := dividend / divisor

			if c.HasFlag("remainder") {
				remainder := math.Mod(dividend, divisor)
				fmt.Printf("%.2f ÷ %.2f = %.2f (remainder: %.2f)\n", dividend, divisor, quotient, remainder)
			} else {
				fmt.Printf("%.2f ÷ %.2f = %.2f\n", dividend, divisor, quotient)
			}

			return nil
		})
	})

	// Power calculation
	app.Command("power", func(cmd *gossie.Command) {
		cmd.Description("Calculate power (base^exponent)")
		cmd.Arg("base").Required().Description("Base number")
		cmd.Arg("exponent").Required().Description("Exponent")
		cmd.Action(func(c *gossie.Context) error {
			base, err := strconv.ParseFloat(c.GetArg("base"), 64)
			if err != nil {
				return fmt.Errorf("invalid base: %s", c.GetArg("base"))
			}

			exponent, err := strconv.ParseFloat(c.GetArg("exponent"), 64)
			if err != nil {
				return fmt.Errorf("invalid exponent: %s", c.GetArg("exponent"))
			}

			result := math.Pow(base, exponent)
			fmt.Printf("%.2f^%.2f = %.2f\n", base, exponent, result)
			return nil
		})
	})

	// Square root
	app.Command("sqrt", func(cmd *gossie.Command) {
		cmd.Description("Calculate square root")
		cmd.Arg("number").Required().Description("Number to find square root of")
		cmd.Action(func(c *gossie.Context) error {
			num, err := strconv.ParseFloat(c.GetArg("number"), 64)
			if err != nil {
				return fmt.Errorf("invalid number: %s", c.GetArg("number"))
			}

			if num < 0 {
				return fmt.Errorf("cannot calculate square root of negative number")
			}

			result := math.Sqrt(num)
			fmt.Printf("√%.2f = %.2f\n", num, result)
			return nil
		})
	})

	// Trigonometric functions
	app.Command("sin", func(cmd *gossie.Command) {
		cmd.Description("Calculate sine of angle in degrees")
		cmd.Arg("angle").Required().Description("Angle in degrees")
		cmd.Action(func(c *gossie.Context) error {
			angle, err := strconv.ParseFloat(c.GetArg("angle"), 64)
			if err != nil {
				return fmt.Errorf("invalid angle: %s", c.GetArg("angle"))
			}

			result := math.Sin(angle * math.Pi / 180)
			fmt.Printf("sin(%.2f°) = %.4f\n", angle, result)
			return nil
		})
	})

	app.Command("cos", func(cmd *gossie.Command) {
		cmd.Description("Calculate cosine of angle in degrees")
		cmd.Arg("angle").Required().Description("Angle in degrees")
		cmd.Action(func(c *gossie.Context) error {
			angle, err := strconv.ParseFloat(c.GetArg("angle"), 64)
			if err != nil {
				return fmt.Errorf("invalid angle: %s", c.GetArg("angle"))
			}

			result := math.Cos(angle * math.Pi / 180)
			fmt.Printf("cos(%.2f°) = %.4f\n", angle, result)
			return nil
		})
	})

	app.Command("tan", func(cmd *gossie.Command) {
		cmd.Description("Calculate tangent of angle in degrees")
		cmd.Arg("angle").Required().Description("Angle in degrees")
		cmd.Action(func(c *gossie.Context) error {
			angle, err := strconv.ParseFloat(c.GetArg("angle"), 64)
			if err != nil {
				return fmt.Errorf("invalid angle: %s", c.GetArg("angle"))
			}

			result := math.Tan(angle * math.Pi / 180)
			fmt.Printf("tan(%.2f°) = %.4f\n", angle, result)
			return nil
		})
	})

	// Statistics commands
	app.Command("stats", func(cmd *gossie.Command) {
		cmd.Description("Statistical calculations")

		cmd.Command("mean", func(subcmd *gossie.Command) {
			subcmd.Description("Calculate arithmetic mean")
			subcmd.Arg("numbers").Multiple().Description("Numbers to average")
			subcmd.Action(func(c *gossie.Context) error {
				args := c.Args()
				if len(args) == 0 {
					return fmt.Errorf("at least one number required")
				}

				var sum float64
				for _, arg := range args {
					num, err := strconv.ParseFloat(arg, 64)
					if err != nil {
						return fmt.Errorf("invalid number: %s", arg)
					}
					sum += num
				}

				mean := sum / float64(len(args))
				fmt.Printf("Mean of %v = %.2f\n", args, mean)
				return nil
			})
		})

		cmd.Command("median", func(subcmd *gossie.Command) {
			subcmd.Description("Calculate median")
			subcmd.Arg("numbers").Multiple().Description("Numbers to find median of")
			subcmd.Action(func(c *gossie.Context) error {
				args := c.Args()
				if len(args) == 0 {
					return fmt.Errorf("at least one number required")
				}

				var numbers []float64
				for _, arg := range args {
					num, err := strconv.ParseFloat(arg, 64)
					if err != nil {
						return fmt.Errorf("invalid number: %s", arg)
					}
					numbers = append(numbers, num)
				}

				// Simple sort for median calculation
				for i := 0; i < len(numbers)-1; i++ {
					for j := i + 1; j < len(numbers); j++ {
						if numbers[i] > numbers[j] {
							numbers[i], numbers[j] = numbers[j], numbers[i]
						}
					}
				}

				var median float64
				n := len(numbers)
				if n%2 == 0 {
					median = (numbers[n/2-1] + numbers[n/2]) / 2
				} else {
					median = numbers[n/2]
				}

				fmt.Printf("Median of %v = %.2f\n", args, median)
				return nil
			})
		})
	})

	app.Run()
}
