package product

type productService struct {
	repo *productRepository
}

func NewProductService(repo *productRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetAllProducts() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetProductByID(id int) (Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) CreateProduct(product Product) (Product, error) {
	return s.repo.Create(product)
}

func (s *productService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
