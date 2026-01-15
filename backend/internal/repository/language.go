package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type LanguageRepository struct {
	*BaseRepository[domain.Language]
}

func NewLanguageRepository(db *gorm.DB) *LanguageRepository {
	return &LanguageRepository{
		BaseRepository: NewBaseRepository[domain.Language](db),
	}
}

func (r *LanguageRepository) FindByCode(ctx context.Context, code string) (*domain.Language, error) {
	var lang domain.Language
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&lang).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

func (r *LanguageRepository) FindActive(ctx context.Context) ([]domain.Language, error) {
	var languages []domain.Language
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).Find(&languages).Error
	return languages, err
}

func (r *LanguageRepository) FindDefault(ctx context.Context) (*domain.Language, error) {
	var lang domain.Language
	err := r.DB.WithContext(ctx).Where("is_default = ?", true).First(&lang).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

type TranslationRepository struct {
	*BaseRepository[domain.Translation]
}

func NewTranslationRepository(db *gorm.DB) *TranslationRepository {
	return &TranslationRepository{
		BaseRepository: NewBaseRepository[domain.Translation](db),
	}
}

func (r *TranslationRepository) FindByLanguageCode(ctx context.Context, langCode string) ([]domain.Translation, error) {
	var translations []domain.Translation
	err := r.DB.WithContext(ctx).Where("language_code = ?", langCode).Find(&translations).Error
	return translations, err
}

func (r *TranslationRepository) FindByModule(ctx context.Context, module string) ([]domain.Translation, error) {
	var translations []domain.Translation
	err := r.DB.WithContext(ctx).Where("module = ?", module).Find(&translations).Error
	return translations, err
}

func (r *TranslationRepository) FindByLanguageAndModule(ctx context.Context, langCode, module string) ([]domain.Translation, error) {
	var translations []domain.Translation
	err := r.DB.WithContext(ctx).
		Where("language_code = ? AND module = ?", langCode, module).
		Find(&translations).Error
	return translations, err
}

func (r *TranslationRepository) FindByKey(ctx context.Context, key string) ([]domain.Translation, error) {
	var translations []domain.Translation
	err := r.DB.WithContext(ctx).Where("key = ?", key).Find(&translations).Error
	return translations, err
}

func (r *TranslationRepository) GetValue(ctx context.Context, langCode, key string) (string, error) {
	var translation domain.Translation
	err := r.DB.WithContext(ctx).
		Where("language_code = ? AND key = ?", langCode, key).
		First(&translation).Error
	if err != nil {
		return "", err
	}
	return translation.Value, nil
}

func (r *TranslationRepository) GetAllKeys(ctx context.Context) ([]string, error) {
	var keys []string
	err := r.DB.WithContext(ctx).Model(&domain.Translation{}).
		Distinct("key").
		Pluck("key", &keys).Error
	return keys, err
}

func (r *TranslationRepository) Upsert(ctx context.Context, translation *domain.Translation) error {
	return r.DB.WithContext(ctx).
		Where("language_code = ? AND key = ?", translation.LanguageCode, translation.Key).
		Assign(map[string]interface{}{
			"value":  translation.Value,
			"module": translation.Module,
		}).
		FirstOrCreate(translation).Error
}

func (r *TranslationRepository) BulkUpsert(ctx context.Context, translations []domain.Translation) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, t := range translations {
			if err := tx.Where("language_code = ? AND key = ?", t.LanguageCode, t.Key).
				Assign(map[string]interface{}{
					"value":  t.Value,
					"module": t.Module,
				}).
				FirstOrCreate(&t).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *TranslationRepository) GetTranslationsMap(ctx context.Context, langCode string) (map[string]string, error) {
	translations, err := r.FindByLanguageCode(ctx, langCode)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, t := range translations {
		result[t.Key] = t.Value
	}
	return result, nil
}
