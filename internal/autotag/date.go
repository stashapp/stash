package autotag

import (
	"context"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

// 다양한 날짜 패턴을 인식하는 정규식들
var datePatterns = []*regexp.Regexp{
	// YYYY-MM-DD 형식 (2020-11-12)
	regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`),
	// YYYYMMDD 형식 (20201112)
	regexp.MustCompile(`(\d{4})(\d{2})(\d{2})`),
	// DD.MM.YYYY 형식 (12.11.2020)
	regexp.MustCompile(`(\d{1,2})\.(\d{1,2})\.(\d{4})`),
}

// ExtractDateFromPath는 파일 경로에서 날짜를 추출합니다.
func ExtractDateFromPath(path string) *time.Time {
	// 파일명만 추출
	filename := filepath.Base(path)
	
	// 각 패턴에 대해 시도
	for i, pattern := range datePatterns {
		matches := pattern.FindStringSubmatch(filename)
		if len(matches) >= 4 {
			var year, month, day int
			var err error
			
			// 패턴에 따라 날짜 파싱 로직 구현
			switch i {
			case 0, 1: // YYYY-MM-DD 또는 YYYYMMDD
				year, err = strconv.Atoi(matches[1])
				if err != nil {
					continue
				}
				month, err = strconv.Atoi(matches[2])
				if err != nil {
					continue
				}
				day, err = strconv.Atoi(matches[3])
				if err != nil {
					continue
				}
			case 2: // DD.MM.YYYY
				day, err = strconv.Atoi(matches[1])
				if err != nil {
					continue
				}
				month, err = strconv.Atoi(matches[2])
				if err != nil {
					continue
				}
				year, err = strconv.Atoi(matches[3])
				if err != nil {
					continue
				}
			}
			
			// 날짜 유효성 검사
			if year < 1900 || year > 2100 || month < 1 || month > 12 || day < 1 || day > 31 {
				continue
			}
			
			// 날짜 객체 생성
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			return &date
		}
	}
	
	return nil
}

// SceneDate는 장면의 파일 경로에서 날짜를 추출하여 설정합니다.
func SceneDate(ctx context.Context, s *models.Scene, rw models.SceneUpdater) error {
	// 이미 날짜가 설정되어 있으면 건너뜀
	if s.Date != nil {
		return nil
	}
	
	// 파일 경로에서 날짜 추출
	date := ExtractDateFromPath(s.Path)
	if date == nil {
		return nil // 날짜를 찾지 못함
	}
	
	// 장면 객체 업데이트
	partial := models.NewScenePartial()
	
	// time.Time을 models.Date로 변환
	dateModel := models.Date{Time: *date}
	partial.Date = models.NewOptionalDate(dateModel)
	
	// 데이터베이스 업데이트
	_, err := rw.UpdatePartial(ctx, s.ID, partial)
	return err
}
