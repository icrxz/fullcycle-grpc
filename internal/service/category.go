package service

import (
	"context"
	"io"

	"github.com/icrxz/fullcycle-grpc/internal/database"
	"github.com/icrxz/fullcycle-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, input *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.Create(input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoriesList{}

	for {
		input, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(categories)
			}
			return err
		}

		category, err := s.CategoryDB.Create(input.Name, input.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
}

func (s *CategoryService) CreateCategoryStreamBidiretional(stream pb.CategoryService_CreateCategoryStreamBidiretionalServer) error {
	for {
		input, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		category, err := s.CategoryDB.Create(input.Name, input.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
		if err != nil {
			return err
		}
	}
}

func (s *CategoryService) GetCategory(ctx context.Context, input *pb.GetCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.GetByID(input.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, input *pb.Blank) (*pb.CategoriesList, error) {
	categories, err := s.CategoryDB.ListAll()
	if err != nil {
		return nil, err
	}

	var categoriesList pb.CategoriesList
	for _, category := range categories {
		categoriesList.Categories = append(categoriesList.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &categoriesList, nil
}
