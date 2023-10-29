package helpers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/mutation"
)

/*
Add dummy data to Dynamo database

@param - table DynamoDB client table

@returns - error If found
*/
func PopulateStoryDatabase(table mutation.TableBasics) error {
	var err error
	// Log actions
	fmt.Println("Populating database with stories.")

	// Upload stories
	err = table.AddStory(models.DynamoStoryDatabaseStruct{
		Id:          GenerateStringUUID(),
		Title:       "descriptive title",
		Description: "descriptive description text",
		Content: &models.StoryContent{
			Id:          GenerateStringUUID(),
			Name:        "Nam blandit magna vel lacinia",
			Description: "Quisque blandit magna vel lacinia fringilla. Mauris sit",
			Children: &[]models.StoryContent{
				{
					Id:          GenerateStringUUID(),
					Name:        "Porttitor quis ultrices tortor",
					Description: "Quisque blandit magna vel lacinia fringilla. Mauris sit amet gravida tellus. Ut sagittis convallis bibendum.",
					Children: &[]models.StoryContent{
						{
							Id:          GenerateStringUUID(),
							Name:        "Nam blandit magna vel lacinia",
							Description: "Let it be known",
							Children:    nil,
						},
						{
							Id:          GenerateStringUUID(),
							Name:        "Euismod amet sapien malesuada",
							Description: "Maecenas lacinia quam eu quam varius semper. Nullam fringilla dapibus ligula, eget porttitor nibh vulputate ut. In hac habitasse platea dictumst. Sed lectus metus, lobortis a ultrices non, malesuada et mauris. Etiam ut facilisis sapien. Praesent iaculis rutrum arcu, at dapibus arcu venenatis a. Mauris ut velit vitae magna commodo convallis ac nec nibh.",
							Children: &[]models.StoryContent{
								{
									Id:          GenerateStringUUID(),
									Name:        "Ullamcorper pulvinar libero",
									Description: "Maecenas lacinia quam eu quam varius semper. Nullam fringilla dapibus ligula, eget porttitor nibh vulputate ut. In hac habitasse platea dictumst. Sed lectus metus, lobortis a ultrices non, malesuada et mauris. Etiam ut facilisis sapien. Praesent iaculis rutrum arcu, at dapibus arcu venenatis a. Mauris ut velit vitae magna commodo convallis ac nec nibh.",
									Children:    nil,
								},
							},
						},
						{
							Id:   GenerateStringUUID(),
							Name: "Fake API",
							Children: &[]models.StoryContent{
								{
									Id:          GenerateStringUUID(),
									Name:        "Nam blandit magna vel lacinia",
									Description: "Etiam eu sollicitudin nisi. Nunc condimentum vel arcu vel sagittis. Maecenas vestibulum volutpat ultricies. Nunc eget purus sapien. Nam sollicitudin nisi sit amet finibus euismod. Suspendisse pretium sapien sit amet mauris vestibulum porttitor. Vivamus vitae purus porttitor, ultrices orci pretium, fringilla orci. Proin facilisis rhoncus mi, eget ullamcorper nibh. Vestibulum condimentum mauris sit amet enim tincidunt, nec vestibulum metus vulputate. Phasellus dui nibh, consequat ut risus ac, facilisis feugiat felis. Donec fermentum, diam in sollicitudin rhoncus, velit arcu volutpat leo, quis lacinia elit metus vitae orci.",
									Children:    nil,
								},
								{
									Id:          GenerateStringUUID(),
									Name:        "Porttitor quis ultrices tortor",
									Description: "Nunc fringilla libero in metus pharetra, a ultrices ipsum pretium. Aliquam hendrerit ex eget risus posuere faucibus. Cras tristique, mauris id vestibulum pulvinar, justo metus luctus urna, id pellentesque mi ligula quis nulla. Fusce ac est justo. Cras eget tempor lectus. Aenean bibendum purus egestas egestas efficitur. Praesent eget tortor non turpis euismod euismod. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Maecenas porttitor vulputate risus et rutrum. Phasellus nec elit lobortis, mollis elit vel, efficitur est. Mauris et pellentesque magna, vitae fermentum urna. Integer tincidunt magna dolor, vitae posuere nunc iaculis et. Aliquam ut sem eu magna gravida imperdiet. Proin sit amet nunc lectus. Duis tristique vulputate elementum.",
									Children:    nil,
								},
							},
						},
						{
							Id:          GenerateStringUUID(),
							Name:        "Quisque",
							Description: "Maecenas lacinia quam eu quam varius semper. Nullam fringilla dapibus ligula, eget porttitor nibh vulputate ut. In hac habitasse platea dictumst. Sed lectus metus, lobortis a ultrices non, malesuada et mauris. Etiam ut facilisis sapien. Praesent iaculis rutrum arcu, at dapibus arcu venenatis a. Mauris ut velit vitae magna commodo convallis ac nec nibh.",
							Children:    nil,
						},
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	err = table.AddStory(models.DynamoStoryDatabaseStruct{
		Id:          GenerateStringUUID(),
		Title:       "Porttitor quis ultrices tortor",
		Description: "Nullam non tempor nisi, ut porta ex. Aenean non mi et nibh feugiat congue id et lacus.",
		Content:     nil})

	if err != nil {
		return err
	}

	err = table.AddStory(models.DynamoStoryDatabaseStruct{
		Id:          GenerateStringUUID(),
		Title:       "website request title",
		Description: "website request description",
		Content: &models.StoryContent{
			Id:          GenerateStringUUID(),
			Name:        "Nam blandit magna vel lacinia",
			Description: "Quisque blandit magna vel lacinia fringilla. Mauris sit",
			Children: &[]models.StoryContent{
				{
					Id:          GenerateStringUUID(),
					Name:        "Porttitor quis ultrices tortor",
					Description: "Quisque blandit magna vel lacinia fringilla. Mauris sit amet gravida tellus. Ut sagittis convallis bibendum.",
					Children: &[]models.StoryContent{
						{
							Id:          GenerateStringUUID(),
							Name:        "Nam blandit magna vel lacinia",
							Description: "Let it be known",
							Children:    nil,
						},
						{
							Id:          GenerateStringUUID(),
							Name:        "Euismod amet sapien malesuada",
							Description: "Maecenas lacinia quam eu quam varius semper. Nullam fringilla dapibus ligula, eget porttitor nibh vulputate ut. In hac habitasse platea dictumst. Sed lectus metus, lobortis a ultrices non, malesuada et mauris. Etiam ut facilisis sapien. Praesent iaculis rutrum arcu, at dapibus arcu venenatis a. Mauris ut velit vitae magna commodo convallis ac nec nibh.",
							Children: &[]models.StoryContent{
								{
									Id:          GenerateStringUUID(),
									Name:        "Ullamcorper pulvinar libero",
									Description: "Maecenas lacinia quam eu quam varius semper. Nullam fringilla dapibus ligula, eget porttitor nibh vulputate ut. In hac habitasse platea dictumst. Sed lectus metus, lobortis a ultrices non, malesuada et mauris. Etiam ut facilisis sapien. Praesent iaculis rutrum arcu, at dapibus arcu venenatis a. Mauris ut velit vitae magna commodo convallis ac nec nibh.",
									Children:    nil,
								},
							},
						},
						{
							Id:   GenerateStringUUID(),
							Name: "Fake API",
							Children: &[]models.StoryContent{
								{
									Id:          GenerateStringUUID(),
									Name:        "Nam blandit magna vel lacinia",
									Description: "Etiam eu sollicitudin nisi. Nunc condimentum vel arcu vel sagittis. Maecenas vestibulum volutpat ultricies. Nunc eget purus sapien. Nam sollicitudin nisi sit amet finibus euismod. Suspendisse pretium sapien sit amet mauris vestibulum porttitor. Vivamus vitae purus porttitor, ultrices orci pretium, fringilla orci. Proin facilisis rhoncus mi, eget ullamcorper nibh. Vestibulum condimentum mauris sit amet enim tincidunt, nec vestibulum metus vulputate. Phasellus dui nibh, consequat ut risus ac, facilisis feugiat felis. Donec fermentum, diam in sollicitudin rhoncus, velit arcu volutpat leo, quis lacinia elit metus vitae orci.",
									Children:    nil,
								},
								{
									Id:          GenerateStringUUID(),
									Name:        "Porttitor quis ultrices tortor",
									Description: "Nunc fringilla libero in metus pharetra, a ultrices ipsum pretium. Aliquam hendrerit ex eget risus posuere faucibus. Cras tristique, mauris id vestibulum pulvinar, justo metus luctus urna, id pellentesque mi ligula quis nulla. Fusce ac est justo. Cras eget tempor lectus. Aenean bibendum purus egestas egestas efficitur. Praesent eget tortor non turpis euismod euismod. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Maecenas porttitor vulputate risus et rutrum. Phasellus nec elit lobortis, mollis elit vel, efficitur est. Mauris et pellentesque magna, vitae fermentum urna. Integer tincidunt magna dolor, vitae posuere nunc iaculis et. Aliquam ut sem eu magna gravida imperdiet. Proin sit amet nunc lectus. Duis tristique vulputate elementum.",
									Children:    nil,
								},
							},
						},
						{
							Id:          GenerateStringUUID(),
							Name:        "Quisque",
							Description: "Maecenas lacinia quam eu quam varius semper. Nullam fringilla dapibus ligula, eget porttitor nibh vulputate ut. In hac habitasse platea dictumst. Sed lectus metus, lobortis a ultrices non, malesuada et mauris. Etiam ut facilisis sapien. Praesent iaculis rutrum arcu, at dapibus arcu venenatis a. Mauris ut velit vitae magna commodo convallis ac nec nibh.",
							Children:    nil,
						},
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func waitForTable(ctx context.Context, db *dynamodb.Client, tn string) error {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tn),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return errors.New("timed out while waiting for table to become active")
	}

	return err
}
