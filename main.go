package main

import (
	"errors"
	"os"

	"github.com/dan-v/rattlesnakeos-stack/stack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version string
var name, region, device, ami, sshKey, spotPrice string
var remove, preventShutdown bool

var rootCmd = &cobra.Command{
	Use:   "rattlesnakeos-stack",
	Short: "Setup AWS infrastructure to build RattlesnakeOS with OTA updates",
	Args: func(cmd *cobra.Command, args []string) error {
		if device != "marlin" && device != "sailfish" && device != "taimen" && device != "walleye" {
			return errors.New("Must specify either marlin|sailfish|taimen|walleye for device type")
		}
		return nil
	},
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		if !remove {
			if err := stack.AWSApply(
				stack.StackConfig{
					Name:            name,
					Region:          region,
					Device:          device,
					AMI:             ami,
					SSHKey:          sshKey,
					SpotPrice:       spotPrice,
					PreventShutdown: preventShutdown,
				},
			); err != nil {
				log.Fatal(err)
			}
		} else {
			stack.AWSDestroy(
				stack.StackConfig{
					Name:   name,
					Region: region,
				},
			)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "name for stack. note: this must be a valid/unique S3 bucket name.")
	rootCmd.MarkFlagRequired("name")
	rootCmd.Flags().StringVarP(&region, "region", "r", "", "aws region for deployment (e.g. us-west-2)")
	rootCmd.MarkFlagRequired("region")
	rootCmd.Flags().StringVarP(&device, "device", "d", "", "device you want to build for: 'marlin' (Pixel XL), 'sailfish' (Pixel), 'taimen' (Pixel 2 XL), 'walleye' (Pixel 2)")
	rootCmd.MarkFlagRequired("device")
	rootCmd.Flags().StringVar(&sshKey, "ssh-key", "", "aws ssh key to add to ec2 spot instances. this is optional but is useful for debugging build issues on the instance.")
	rootCmd.Flags().StringVar(&spotPrice, "spot-price", ".80", "spot price for build ec2 instances. if this value is too low you may not obtain an instance or it may terminate during a build.")
	rootCmd.Flags().StringVar(&ami, "ami", "", "ami id to use for build environment. this is optional as correct ubuntu ami for region will be chosen by default.")
	rootCmd.Flags().BoolVar(&remove, "remove", false, "cleanup/destroy all deployed aws resources.")
	rootCmd.Flags().BoolVar(&preventShutdown, "prevent-shutdown", false, "for debugging purposes only - will prevent ec2 instance from shutting down after build.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
